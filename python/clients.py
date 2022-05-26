#!/usr/bin/env python

import asyncio
from dataclasses import dataclass
import datetime

import json
import time

import websockets


@dataclass
class pixel:
    x: int
    y: int
    color: int
    timestamp: int
    userid: int


async def main():
    message = pixel(
        x=0,
        y=1,
        color=0,
        timestamp=int(time.time()),
        userid=0,
    )
    async with websockets.connect("ws://localhost:8080") as websocket:
        print(message)
        await websocket.send(json.dumps(message.__dict__))
        await websocket.recv()


if __name__ == "__main__":
    asyncio.get_event_loop().run_until_complete(main())
