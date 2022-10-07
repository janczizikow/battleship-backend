X-Token: device-id

{
  code: asdfasdf-asdasdf-asdfas
}


1. POST /rooms


payload:
```
{
  name: string (optional, if not provided generates a random one)
  player1Id: '',
}
```

response:

```
{
  name: '',
  code: 1234
}
```

```
{
  errors: [
    {player1Id: 'is required'}
  ]
}
```

2. POST /rooms/:roomCode/join


request:

```
{
  player2Id: '',
}
```

response:
```
{
  name: '',
  player1Id: '',
  player2Id: '',
}
```

404:

{
  error: 'room not found'
}