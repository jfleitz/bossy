# bossy
GO Source for Game Plan Mike Bossy pinball machine utilizing the goflip library.

## Introduction
The Mike Bossy pinball machine never made it to production. There was one made for the 1982 AMOA show (confirmed - playmeter magazine Feb 15,1982 ed), however it doesn't match what the flyer said the game was going to be.

The original roms cannot be found. Therefore, the game is being redone with a more hockey feel, without making any changes to the playfield design.

Check out the playfield layout on ipdb: http://www.ipdb.org/machine.cgi?id=1596


## Overview of Rules
* Pucks - The playfield has white spots. These spots are considered to be hockey pucks. Only one is lit at a time. When you "hit the puck", the next MIKE BOSSY letter is lit, and another puck location is chosen random.

* 9 second warm up period - This is from the original game. Basically a ball save for 9 seconds. You can't hit a goal at this time (or it won't count). It allows you to build up the number of pucks collected (Mike Bossy letters)

* Goal - When a goal is hit, the pucks collected (MIKE BOSSY letters) are cashed in for points.

* LILCO Line Substition - When this mode is on, the pinball must be hit to one of the left inlane flipper lanes, one of the right inlane flipper lanes, and the saucer in the top right of the playfield. This refers to when Bryan Trottier, Clark Gillies, and Mike Bossy played at the same time forming the LILCO (Long Island Lighting Company) line (name was given since the goal light stayed on all of the time).
