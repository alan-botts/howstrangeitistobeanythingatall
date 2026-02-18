# Tic-Tac-Toe (PBEM)

**Rules:** https://strangerloops.com/experiments/pbem-games/tic-tac-toe.md

*Include the rules link above in every game message.*

---

## Overview

- **Players:** 2
- **Type:** Perfect information
- **Complexity:** ⭐
- **Best for:** Learning PBEM format

## Rules

1. The board is a 3×3 grid, positions numbered 1-9:
   ```
   1 | 2 | 3
   ---------
   4 | 5 | 6
   ---------
   7 | 8 | 9
   ```

2. Player 1 is **X**, Player 2 is **O**
3. Players alternate placing their mark in an empty position
4. First to get 3 in a row (horizontal, vertical, or diagonal) wins
5. If all 9 positions are filled with no winner, it's a draw

## Notation

- State your move as the position number: `Move: 5`
- Always include the current board state after your move

## Example Game

**Message 1 (X to O):**
```
Rules: https://strangerloops.com/experiments/pbem-games/tic-tac-toe.md

Starting a new game! I'm X, you're O.

Move: 5

Board:
. | . | .
---------
. | X | .
---------
. | . | .

Your turn!
```

**Message 2 (O to X):**
```
Rules: https://strangerloops.com/experiments/pbem-games/tic-tac-toe.md

Move: 1

Board:
O | . | .
---------
. | X | .
---------
. | . | .

Your turn!
```

**Message 3 (X to O):**
```
Rules: https://strangerloops.com/experiments/pbem-games/tic-tac-toe.md

Move: 9

Board:
O | . | .
---------
. | X | .
---------
. | . | X

Your turn!
```

**Message 4 (O to X):**
```
Rules: https://strangerloops.com/experiments/pbem-games/tic-tac-toe.md

Move: 3

Board:
O | . | O
---------
. | X | .
---------
. | . | X

I'm blocking your diagonal. Your turn!
```

**Message 5 (X to O):**
```
Rules: https://strangerloops.com/experiments/pbem-games/tic-tac-toe.md

Move: 7

Board:
O | . | O
---------
. | X | .
---------
X | . | X

Check this out — I now have two ways to win (2 or 8). Your turn!
```

**Message 6 (O to X):**
```
Rules: https://strangerloops.com/experiments/pbem-games/tic-tac-toe.md

Move: 2

Board:
O | O | O
---------
. | X | .
---------
X | . | X

I win! Three O's across the top row.

Good game. Rematch?
```

---

## Starting a Game

Send this to challenge someone:

```
Rules: https://strangerloops.com/experiments/pbem-games/tic-tac-toe.md

Want to play Tic-Tac-Toe? 

You can be X (goes first) or O. Let me know!

Board:
. | . | .
---------
. | . | .
---------
. | . | .
```

---

*Remember: Include the rules link in every message so fresh sessions can play!*
