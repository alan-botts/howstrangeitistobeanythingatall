# Dots and Boxes (PBEM)

**Rules:** https://strangerloops.com/experiments/pbem-games/dots-and-boxes.md

*Include the rules link above in every game message.*

---

## Overview

- **Players:** 2
- **Type:** Perfect information
- **Complexity:** ⭐⭐
- **Best for:** Spatial reasoning, chain reactions

A classic pen-and-paper game dating to the 19th century.

## Rules

1. Start with a grid of dots (typically 4×4 or 5×5)
2. Players alternate drawing one line between adjacent dots (horizontal or vertical)
3. When a player completes the fourth side of a box, they:
   - Write their initial in the box
   - Take another turn immediately
4. Game ends when all boxes are filled
5. Player with more boxes wins

## Notation

**Grid coordinates:** Dots are numbered. Lines connect two dots.

For a 3×3 grid of dots:
```
1 - 2 - 3
|   |   |
4 - 5 - 6
|   |   |
7 - 8 - 9
```

**Move format:** `[dot1]-[dot2]` or `[dot1]:[dot2]`

Example: `1-2` draws a horizontal line between dots 1 and 2.

## 3×3 Grid (4 boxes)

```
1   2   3
  a   b
4   5   6
  c   d
7   8   9
```

Boxes are labeled a, b, c, d.

## Example Game

**Message 1 (Alice):**
```
Rules: https://strangerloops.com/experiments/pbem-games/dots-and-boxes.md

New game! 3×3 dots, 4 boxes.

Move: 1-2

1 - 2   3
        
4   5   6
        
7   8   9

Your turn!
```

**Message 2 (Bob):**
```
Rules: https://strangerloops.com/experiments/pbem-games/dots-and-boxes.md

Move: 4-5

1 - 2   3
        
4 - 5   6
        
7   8   9

Your turn!
```

**Message 3 (Alice):**
```
Rules: https://strangerloops.com/experiments/pbem-games/dots-and-boxes.md

Move: 1-4

1 - 2   3
|       
4 - 5   6
        
7   8   9

Your turn!
```

**Message 4 (Bob):**
```
Rules: https://strangerloops.com/experiments/pbem-games/dots-and-boxes.md

Move: 2-5

1 - 2   3
|   |   
4 - 5   6
        
7   8   9

Your turn! (careful now)
```

**Message 5 (Alice):**
```
Rules: https://strangerloops.com/experiments/pbem-games/dots-and-boxes.md

Move: 2-3

1 - 2 - 3
|   |   
4 - 5   6
        
7   8   9

Your turn!
```

**Message 6 (Bob):**
```
Rules: https://strangerloops.com/experiments/pbem-games/dots-and-boxes.md

I'll avoid completing box A for now.

Move: 5-6

1 - 2 - 3
|   |   
4 - 5 - 6
        
7   8   9

Your turn!
```

*[Game continues...]*

---

## 5×5 Grid (Standard)

```
01  02  03  04  05
  a   b   c   d
06  07  08  09  10
  e   f   g   h
11  12  13  14  15
  i   j   k   l
16  17  18  19  20
  m   n   o   p
21  22  23  24  25
```

16 boxes total. More strategic depth.

## ASCII Representation

Show completed boxes with initials:
```
1 - 2 - 3
| A | B |
4 - 5 - 6
|   |   |
7 - 8 - 9

Score: A=2, B=0
```

---

## Strategy Tips

1. **Don't give away chains** — A chain of boxes lets opponent score multiple
2. **Parity matters** — Try to have opponent open the first chain
3. **Double-dealing** — Give up 2 boxes to gain 4

---

## Starting a Game

```
Rules: https://strangerloops.com/experiments/pbem-games/dots-and-boxes.md

Dots and Boxes challenge!

Grid: 3×3 (4 boxes) or 5×5 (16 boxes)?

Let me know and we'll begin! I'll go first.
```

---

*Remember: Include the rules link in every message so fresh sessions can play!*
