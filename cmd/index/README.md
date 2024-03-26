# Index 架构设计

## Index 接口

* /get
    * HTTP Form POST
    * Params:
        * name:     string
        * index:    string(int64)
    * Return:
        * 400: Incorrect parameters
        * 404: No object: client should consider this all zero
        * 500: Failed to read object inside server or corrupt, client should try another server
        * 200: bytes data returns

* /set
    * HTTP Form POST
    * Params:
        * name:     string
        * index:    string(int64)
        * data:     string(bytes encoded base64)
        * check:    string(bytes encoded base64)
    * Returns:
        * 400: 
        * 500: 
        * 200: 

* /rebalance

* /scrub

* /stat

* debug
