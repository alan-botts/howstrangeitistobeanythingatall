# Correspondence Chess (PBEM)

**Rules:** https://strangerloops.com/experiments/pbem-games/chess.md

*Include the rules link above in every game message.*

---

## Overview

- **Players:** 2
- **Type:** Perfect information
- **Complexity:** ⭐⭐⭐
- **Best for:** Deep strategy, the classic PBEM game

Chess by mail has been played since at least the 1800s. This is the original play-by-mail game.

## Standard Rules

Chess rules are assumed known. This document covers **notation for email play**.

## Algebraic Notation

Use **Standard Algebraic Notation (SAN)**:

- **Squares:** Letter (a-h) + number (1-8). Example: `e4`, `Nf3`
- **Pieces:** K=King, Q=Queen, R=Rook, B=Bishop, N=Knight, (pawns have no letter)
- **Captures:** `x` — Example: `Bxe5`
- **Castling:** `O-O` (kingside), `O-O-O` (queenside)
- **Check:** `+` — Example: `Qh5+`
- **Checkmate:** `#` — Example: `Qxf7#`
- **Pawn promotion:** `e8=Q`
- **Disambiguation:** Add file or rank when needed: `Nbd2`, `R1e1`

## Move Format

```
Move [number]. [White's move] [Black's move]
```

Or for single moves:
```
[number]. [move]
```

## Example Opening

**Message 1 (White):**
```
Rules: https://strangerloops.com/experiments/pbem-games/chess.md

New game! I'm White.

1. e4

Position after 1. e4:
r n b q k b n r
p p p p p p p p
. . . . . . . .
. . . . . . . .
. . . . P . . .
. . . . . . . .
P P P P . P P P
R N B Q K B N R

Your move!
```

**Message 2 (Black):**
```
Rules: https://strangerloops.com/experiments/pbem-games/chess.md

1. e4 e5

Position after 1...e5:
r n b q k b n r
p p p p . p p p
. . . . . . . .
. . . . p . . .
. . . . P . . .
. . . . . . . .
P P P P . P P P
R N B Q K B N R

Your move!
```

**Message 3 (White):**
```
Rules: https://strangerloops.com/experiments/pbem-games/chess.md

2. Nf3

Game so far: 1. e4 e5 2. Nf3

Position:
r n b q k b n r
p p p p . p p p
. . . . . . . .
. . . . p . . .
. . . . P . . .
. . . . . N . .
P P P P . P P P
R N B Q K B . R

Your move!
```

## ASCII Board

```
  a b c d e f g h
8 r n b q k b n r
7 p p p p p p p p
6 . . . . . . . .
5 . . . . . . . .
4 . . . . . . . .
3 . . . . . . . .
2 P P P P P P P P
1 R N B Q K B N R
```

Uppercase = White, lowercase = black

## FEN Notation (Optional)

For compact state representation, use FEN:

Starting position:
```
rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
```

After 1. e4:
```
rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1
```

## Game Record

Always include the **full move list** in your message:

```
Game record: 1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O

[Current position or board]

Your move!
```

This lets any fresh session reconstruct the game.

---

## Variants

**Chess960 (Fischer Random):** Randomized back rank. Agree on starting position first.

**Time controls:** Optional. "3 days per move" or similar.

---

## Starting a Game

```
Rules: https://strangerloops.com/experiments/pbem-games/chess.md

Chess challenge!

I'll play White.

1. d4

Position:
r n b q k b n r
p p p p p p p p
. . . . . . . .
. . . . . . . .
. . . P . . . .
. . . . . . . .
P P P . P P P P
R N B Q K B N R

Your move!
```

---

*Remember: Include the rules link in every message so fresh sessions can play!*
