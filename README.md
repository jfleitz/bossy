# bossy
GO Source for Game Plan Mike Bossy pinball machine utilizing the goflip library.

## Introduction
The Mike Bossy pinball machine never made it to production. There was one made for the 1982 AMOA show (confirmed - Roger Sharpe wrote a review of it for  Play Meter magazine Feb 15,1982), however it doesn't match what the flyer said the game was going to be, and the review didn't make it sound like it was "a true Hockey experience".

The original roms cannot be found. There are a few "Bossy game roms" floating out there, but they are garbage. Therefore, the game is being redone with a more hockey feel, without making any changes to the playfield design. 

Check out the playfield layout on ipdb: http://www.ipdb.org/machine.cgi?id=1596

## Overview of Rules
* Pucks - The playfield has white spots. These spots are considered to be hockey pucks. Only one is lit at a time. When you "hit the puck", the next MIKE BOSSY letter is sequentially lit, and another puck location is chosen random. When choosing at random, the "MIKE BOSSY" lights are lit up in a circular fashion until the first target hit after 2 MIKE BOSSY Light sequences(1).
    * (1) This is to prevent someone from capturing the ball and seeing where the lit puck will be.


* 9 second warm up period - This is from the original game. Basically a ball save for 9 seconds. You can't hit a goal at this time (or it won't count). It allows you to build up the number of pucks collected (Mike Bossy letters) when you do hit a goal. This starts after the ball leaves the shooter lane switch (switch up event)

* Goal - When a goal is hit, the pucks collected (MIKE BOSSY letters) are cashed in for points. The MIKE BOSSY letters are cleared back to zero (2). A 5000 Bonus light at the bottom is lit up for each goal scored. Goals are 10,000 for each letter lit plus 10,000 (up to 100,000 points)
    * (2) Even though the letters are cleared, the total "Letter count" for each player is kept during the ball being played, so that the total Letters collected will be awarded in the bonus calculation.

* Bonus  - For each goal you get 5000 bonus light lit. At the end of the ball, the total number of letters * the 5000 will be awarded as a bonus. Example, if you spell MIKE BOSSY twice, and get MIK = that is 21 letters. You have 2 of the 5000 bonus lights lit, so you get 21 times 10,000 = 210,000 for a bonus.

* Overtime awards. On the last ball, the overtime lights are lit. For each time you hit one, 1 second is added to the overtime. If there are more than one player playing, whomever has the higher overtime value gets the timed extra ball.

* Hat Trick - If you hit 3 goals in the same game by the same flipper.
    * For the lower flippers, a goal must be scored within 3 seconds after the ball goes from one of the flipper inlanes to be counted toward the "Hat Trick count" for that flipper.
    * For the upper flipper, a goal must be scored within 2 seconds after the ball is ejected from the saucer to be counted toward the "Hat Trick count" of the upper flipper. 
    * For each Hat Trick (3 shots by the same flipper), 22,000 points is awarded. 

* Trio Grande - This is to pay homage to the best hokcey line in NHL History; Bryan Trottier, Clark Gillies, and Mike Bossy forming the Trio Grande (4 consecutive Stanley Cup wins). Similar to Hat Trick, this is based on timing from the flipper inlane switches and the saucer switch. You must complete each within 5 seconds of each other. Once you do so, the Goal light will blink to signify you have LILCO Line ready. Hitting the goal will award you with 220,000 points.
    * In any order, you must hit within 5 seconds of each other:
        * Either Left Flipper Inlane switch
        * Either Right Flipper Inlane switch
        * The Upper Right saucer
    * This will flash the Goal Light and make the next goal 220,000 points
