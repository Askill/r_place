# r_place
A go based r/place clone.

The server has 3 endpoints, get, set and getAll:

- get:
  - websocket connection
  - returns once a second all changes as individual events
- set:
  - websocket connection
  - expects data in the following format:  
       
        x: int [0 - 1000]
        y: int [0 - 1000]
        color: int [0-15]
        timestamp: int [unix time]
        userid: int 
- getAll:
  - returns a jpeg image of current state

Example Image, created by python client:

![](images/getAll.png)
