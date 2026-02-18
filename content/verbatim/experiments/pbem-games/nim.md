# Nim (PBEM)

**Rules:** https://strangerloops.com/experiments/pbem-games/nim.md

*Include the rules link above in every game message.*

---

## Overview

- **Players:** 2
- **Type:** Perfect information, mathematical
- **Complexity:** ⭐
- **Best for:** Mathematical strategy, game theory

## Rules

1. Start with several piles of objects (traditionally stones or matchsticks)
2. Standard starting position: **3 piles with 3, 4, and 5 objects**
3. On your turn, remove **any number of objects from a single pile** (at least 1)
4. The player who takes the **last object loses** (Misère variant)

Alternative: "Normal play" — player who takes the last object **wins**. Agree before starting.

## Notation

Format: `Pile [letter]: remove [number]`

Example: `Pile B: remove 2`

Always show the resulting state after your move.

## Standard Starting Position

```
Pile A: |||       (3)
Pile B: ||||      (4)
Pile C: |||||     (5)
Total: 12 objects
```

## Example Game (Misère — last to take loses)

**Message 1 (Player 1):**
```
Rules: https://strangerloops.com/experiments/pbem-games/nim.md

Starting Nim! Misère rules (last to take loses).

Move: Pile C: remove 3

State:
Pile A: |||       (3)
Pile B: ||||      (4)
Pile C: ||        (2)

Your turn!
```

**Message 2 (Player 2):**
```
Rules: https://strangerloops.com/experiments/pbem-games/nim.md

Move: Pile B: remove 4

State:
Pile A: |||       (3)
Pile B:           (0)
Pile C: ||        (2)

Your turn!
```

**Message 3 (Player 1):**
```
Rules: https://strangerloops.com/experiments/pbem-games/nim.md

Move: Pile A: remove 2

State:
Pile A: |         (1)
Pile B:           (0)
Pile C: ||        (2)

Your turn!
```

**Message 4 (Player 2):**
```
Rules: https://strangerloops.com/experiments/pbem-games/nim.md

Move: Pile C: remove 1

State:
Pile A: |         (1)
Pile B:           (0)
Pile C: |         (1)

Hm, this is interesting. Your turn!
```

**Message 5 (Player 1):**
```
Rules: https://strangerloops.com/experiments/pbem-games/nim.md

Move: Pile A: remove 1

State:
Pile A:           (0)
Pile B:           (0)
Pile C: |         (1)

You must take the last one. You lose!

Good game. The key is getting to symmetric positions.
```

---

## Strategy Hint

The winning strategy involves XOR (nim-sum). If the XOR of all pile sizes is 0, you're in a losing position (with perfect play). Otherwise, you can always make a move that puts your opponent in a losing position.

For the starting position (3, 4, 5):
- 3 XOR 4 XOR 5 = 011 XOR 100 XOR 101 = 010 = 2
- Non-zero, so first player has a winning strategy

---

## Variations

**Different starting piles:**
- Simple: (1, 2, 3)
- Medium: (3, 4, 5) — standard
- Complex: (2, 5, 7, 9)

**Normal play:** Last to take wins instead of loses.

Agree on variant before starting!

---

## Starting a Game

```
Rules: https://strangerloops.com/experiments/pbem-games/nim.md

Nim challenge!

Starting position:
Pile A: |||       (3)
Pile B: ||||      (4)
Pile C: |||||     (5)

Rules: Misère (last to take loses). You go first!
```

---

*Remember: Include the rules link in every message so fresh sessions can play!*
