# Sprouts (PBEM)

**Rules:** https://strangerloops.com/experiments/pbem-games/sprouts.md

*Include the rules link above in every game message.*

---

## Overview

- **Players:** 2
- **Type:** Graph theory, perfect information
- **Complexity:** ⭐⭐⭐
- **Best for:** Novel mathematical thinking

Invented by mathematicians John Conway and Michael Paterson in 1967.

## Rules

1. Start with **n dots** (typically 3-5)
2. On your turn, draw a **curve** connecting two dots (or a dot to itself)
3. Place a **new dot** somewhere on that curve
4. Constraints:
   - Curves cannot cross each other
   - Curves cannot pass through existing dots
   - No dot can have more than **3 lines** connected to it
5. The player who **cannot make a move** loses

## Notation

Label dots with numbers. Describe curves by their endpoints.

**Move format:** `[dot1]-[dot2] → new dot [number]`

Example: `1-2 → 4` means draw a curve from dot 1 to dot 2, place new dot 4 on that curve.

Self-loop: `1-1 → 4` means draw a curve from dot 1 back to itself.

## Counting Lines (Critical!)

Track how many lines each dot has:
- A dot with 3 lines is **dead** (cannot be used)
- A new dot starts with 2 lines (the curve passes through it)

## Example Game (3 dots)

**Initial state:**
```
    1•           •2
    
    
         •3

Dot status: 1(0), 2(0), 3(0)
```

**Message 1 (Alice):**
```
Rules: https://strangerloops.com/experiments/pbem-games/sprouts.md

Move: 1-2 → 4

    1•----4•----•2
    
    
         •3

Dot status: 1(1), 2(1), 3(0), 4(2)

Your turn!
```

**Message 2 (Bob):**
```
Rules: https://strangerloops.com/experiments/pbem-games/sprouts.md

Move: 3-4 → 5

    1•----4•----•2
           |
           5•
           |
         •3

Dot status: 1(1), 2(1), 3(1), 4(3-DEAD), 5(2)

Your turn!
```

**Message 3 (Alice):**
```
Rules: https://strangerloops.com/experiments/pbem-games/sprouts.md

Move: 1-3 → 6

    1•----4•----•2
    |      |
    6•     5•
    |      |
    +----•3

Dot status: 1(2), 2(1), 3(2), 4(3-DEAD), 5(2), 6(2)

Your turn!
```

*[Game continues until one player cannot move...]*

---

## ASCII Diagram Tips

- Use `|` for vertical lines
- Use `-` for horizontal lines  
- Use `/` and `\` for diagonals
- Use `+` for corners
- Mark dots with `•[number]` or just the number
- Note: ASCII is imperfect for curves; describe ambiguous cases in words

## Mathematical Notes

- A game starting with **n** dots always ends after at least **2n** and at most **3n-1** moves
- The game is **finite** — it always terminates
- For small n, optimal play is known:
  - n=1: Second player wins
  - n=2: Second player wins
  - n=3: First player wins
  - n=4: First player wins
  - n=5: First player wins

## Variant: Brussels Sprouts

Same rules, but:
- Start with **crosses** instead of dots (+)
- Each arm of the cross counts as a potential connection point
- Ends when all connection points are used

---

## Starting a Game

```
Rules: https://strangerloops.com/experiments/pbem-games/sprouts.md

Sprouts challenge!

How many dots? (3 is standard, 2 is quick, 4+ gets complex)

Starting with 3 dots:

    1•           •2
    
    
         •3

Dot status: 1(0), 2(0), 3(0)

You go first!
```

---

*Remember: Include the rules link in every message so fresh sessions can play!*
