package main

import (
	"bufio"
	"fmt"
	"go_utils/utils"
	"log"
	"net"
	"strings"
)

const maxPlayers = 3

var game Game

func init() {
	game = Game{
		players:     make([]*Player, maxPlayers),
		currentTurn: 0,
		landlord:    -1,
		bottomCards: []string{},
		gameState:   GameStateWaiting,
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		utils.LogPrintError("Error starting server: ", err)
		panic(err)
	}
	defer listener.Close()

	utils.LogPrintInfo("Server started on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.LogPrintError("Error accepting connection: ", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	player := &Player{
		conn:   conn,
		name:   "",
		hand:   make([]string, 0),
		status: PlayerStatusWaiting,
	}

	game.mu.Lock()
	playerDidJoinGame := false
	for i, p := range game.players {
		if p == nil || p.status == PlayerStatusDisconnected {
			game.players[i] = player
			player.name = fmt.Sprintf("Player%d", i+1)
			log.Printf("Accepted connection from %s as %s", conn.RemoteAddr(), player.name)
			playerDidJoinGame = true
			break
		}
	}
	if !playerDidJoinGame {
		fmt.Fprintln(conn, "Game already started, please try later.")
		game.mu.Unlock()
		log.Printf("Connection from %s refused, game already started", conn.RemoteAddr())
		return
	}
	game.mu.Unlock()

	fmt.Fprintln(conn, "Welcome", player.name)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			continue
		}

		handlePlayerInput(player, input)
	}

	handlePlayerDisconnection(player)
}

func handlePlayerInput(player *Player, input string) {
	game.mu.Lock()
	defer game.mu.Unlock()

	if game.gameState == GameStateWaiting && len(game.players) == maxPlayers {
		startGame()
	}

	if player.status == PlayerStatusWaiting {
		fmt.Fprintln(player.conn, "Waiting for other players...")
		return
	}

	log.Printf("%s: %s", player.name, input)

	switch game.gameState {
	case GameStateCalling:
		handleCallLandlordInput(player, input)
	case GameStatePlaying:
		handlePlayCardsInput(player, input)
	}
	// Handle input based on game state
	// Implement calling landlord, playing, etc.
}

func handlePlayCardsInput(player *Player, input string) {
	cards := strings.Split(input, " ")
	if !isValidPlay(cards) {
		fmt.Fprintln(player.conn, "Invalid play, try again.")
		return
	}

	playCards(player, cards)
}

func playCards(player *Player, cards []string) {
	game.mu.Lock()
	defer game.mu.Unlock()

	if !removeCardsFromHand(player, cards) {
		fmt.Fprintln(player.conn, "Failed to play cards, invalid cards in hand.")
		return
	}

	game.lastPlayedCards = cards

	log.Printf("%s played: %v", player.name, cards)

	if len(player.hand) == 0 {
		if player == game.players[game.landlord] {
			log.Println("Landlord wins!")
			resetGame()
		} else {
			log.Println("Farmers win!")
			resetGame()
		}
		return
	}

	game.currentTurn = (game.currentTurn + 1) % maxPlayers
	nextPlayer := game.players[game.currentTurn]
	fmt.Fprintf(nextPlayer.conn, "Your turn to play.\n")
}

func removeCardsFromHand(player *Player, cards []string) bool {
	player.mu.Lock()
	defer player.mu.Unlock()

	cardCount := make(map[string]int)
	for _, card := range player.hand {
		cardCount[card]++
	}

	for _, card := range cards {
		if cardCount[card] == 0 {
			return false
		}
		cardCount[card]--
	}

	for _, card := range cards {
		player.hand = removeCard(player.hand, card)
	}

	return true
}

func removeCard(hand []string, card string) []string {
	index := -1
	for i, c := range hand {
		if c == card {
			index = i
			break
		}
	}

	if index != -1 {
		return append(hand[:index], hand[index+1:]...)
	}

	return hand
}

func isValidPlay(rawCards []string) bool {
	cards := make([]Card, len(rawCards))
	for i, rawCard := range rawCards {
		value, exists := CardValueMap[rawCard]
		if !exists {
			return false
		}
		cards[i] = Card{Value: value, Suit: ""}
	}

	switch len(cards) {
	case 1:
		return isSingle(cards)
	case 2:
		return isPair(cards) || isRocket(cards)
	case 3:
		return isThreeOfAKind(cards)
	case 4:
		return isThreeWithOne(cards) || isBomb(cards)
	case 5:
		return isStraight(cards)
	default:
		if isStraight(cards) {
			return true
		}
		if isPlane(cards) {
			return true
		}
		// 添加更多复杂牌型逻辑
		return false
	}
}

func handlePlayerDisconnection(player *Player) {
	game.mu.Lock()
	defer game.mu.Unlock()

	player.status = PlayerStatusDisconnected

	remainingPlayers := 0
	for _, p := range game.players {
		if p != nil && p.status != PlayerStatusDisconnected {
			remainingPlayers++
		}
	}

	if remainingPlayers == 0 {
		log.Println("All players disconnected. Resetting game.")
		resetGame()
		// TODO, this will reset game.mu
	}
}

func handleCallLandlordInput(player *Player, input string) {
	switch input {
	case "y":
		player.mu.Lock()
		player.status = PlayerStatusPlaying
		player.mu.Unlock()
		game.landlord = getPlayerIndex(player)
		fmt.Fprintf(player.conn, "You are the landlord\n")
		game.gameState = GameStatePlaying
		dealBottomCards(game.landlord) // 地主将获得三张底牌
	case "n":
		player.mu.Lock()
		player.status = PlayerStatusRejected
		player.mu.Unlock()
		nextPlayerIndex := (getPlayerIndex(player) + 1) % maxPlayers
		if allPlayersRejected(nextPlayerIndex) {
			log.Println("All players rejected becoming the landlord. Restarting game.")
			resetGame()
			return
		} else {
			nextPlayer := game.players[nextPlayerIndex]
			fmt.Fprintf(nextPlayer.conn, "Do you want to be the landlord (y/n)?\n")
		}
	}
}

func allPlayersRejected(lastIndex int) bool {
	for i := 0; i < maxPlayers; i++ {
		playerIndex := (lastIndex + i) % maxPlayers
		if game.players[playerIndex].status != PlayerStatusRejected {
			return false
		}
	}
	return true
}

func getPlayerIndex(player *Player) int {
	for i, p := range game.players {
		if p == player {
			return i
		}
	}
	return -1
}

func dealBottomCards(landlordIndex int) {
	landlord := game.players[landlordIndex]
	landlord.hand = append(landlord.hand, game.bottomCards...)
}

func startGame() {
	game.gameState = GameStateDealing
	dealCards()

	game.gameState = GameStateCalling
	// Start calling landlord logic
}

func dealCards() {
	cards := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "0", "J", "Q", "K"}
	deck := append(cards, cards...)
	deck = append(deck, cards...)
	deck = append(deck, cards...)
	deck = append(deck, "+", "-")

	// Implement card shuffling and dealing logic
	for i, player := range game.players {
		player.hand = deck[i*17 : (i+1)*17]
	}

	game.bottomCards = deck[51:]
}

func resetGame() {
	game = Game{
		players:     make([]*Player, maxPlayers),
		currentTurn: 0,
		landlord:    -1,
		bottomCards: []string{},
		gameState:   GameStateWaiting,
	}
}
