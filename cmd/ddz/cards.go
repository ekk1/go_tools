package main

import "sort"

type Card struct {
	Value int
	Suit  string
}

var CardValueMap = map[string]int{
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
	"0": 10,
	"J": 11,
	"Q": 12,
	"K": 13,
	"A": 14,
	"2": 15,
	"-": 16,
	"+": 17,
}

func isSingle(cards []Card) bool {
	return len(cards) == 1
}

func isPair(cards []Card) bool {
	return len(cards) == 2 && cards[0].Value == cards[1].Value
}

func isStraight(cards []Card) bool {
	if len(cards) < 5 {
		return false
	}

	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Value < cards[j].Value
	})

	for i := 1; i < len(cards); i++ {
		if cards[i].Value != cards[i-1].Value+1 {
			return false
		}
	}

	return true
}

func isThreeOfAKind(cards []Card) bool {
	return len(cards) == 3 && cards[0].Value == cards[1].Value && cards[0].Value == cards[2].Value
}

func isBomb(cards []Card) bool {
	return len(cards) == 4 && cards[0].Value == cards[1].Value && cards[0].Value == cards[2].Value && cards[0].Value == cards[3].Value
}

func isThreeWithOne(cards []Card) bool {
	if len(cards) != 4 {
		return false
	}
	countMap := make(map[int]int)
	for _, card := range cards {
		countMap[card.Value]++
	}
	return len(countMap) == 2 && (countMap[cards[0].Value] == 3 || countMap[cards[1].Value] == 3)
}

func isPlane(cards []Card) bool {
	if len(cards)%3 != 0 || len(cards) < 6 {
		return false
	}
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Value < cards[j].Value
	})
	countMap := make(map[int]int)
	for _, card := range cards {
		countMap[card.Value]++
	}
	for _, count := range countMap {
		if count != 3 {
			return false
		}
	}
	for i := 1; i < len(cards)/3; i++ {
		if cards[i*3].Value != cards[(i-1)*3].Value+1 {
			return false
		}
	}
	return true
}

func isRocket(cards []Card) bool {
	return len(cards) == 2 && ((cards[0].Value == 16 && cards[1].Value == 17) || (cards[0].Value == 17 && cards[1].Value == 16))
}
