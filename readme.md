Small project to 'harden' a simple opening a file and returning a DISTINCT count of a fake log file. Basically how hard can we go with go in doing a simple Log DISTINCT count.
Meant to showcase some adv tech on working at scale, handling edge cases.


1. Can handle any size file, the way we load is line by line wiht bufio, so a 200gb file would be no problem
2. Outputs to CLI the count of lines processed (total) and any lines that may have been skipped due to corrption or error.
3. Opens 'log' file in read only format to allow it to continue writing.
4. defer file.close() ensures that if we crash midway, file is released. Prevents any 'too many open files' errs on linux
5. write.FLUSH is critical as we want to prevent file empty or cutoffs as GO keeps data in RAM until buffer is full (4kb or so)


Issues
1. MAP bottle neck, if a log contains 50gb of UNIQUE strings, this will be killed by OOM, unless our ram met 50gb needed to hold in memory.
2. Scanner Limit - bufio.Scanner has a max token size(64KB in a line), bound to be triggered by some ungodly JSON. 
Would cause only a error in the current form.


Note that we are currently exporting to CLI, that would be a OOM killer for sure for a 10gb log file.
Depositing the file would be the suprior choice.



Further advances -
1. Sorting could be done at the end via a linux process in bash, just call bash sort uniq -c, if it needs more ram it swaps to disk
2. Use a database, its what its made for. Throw it into Sqlite and this would probally be easy.
3. ADV would be non-exact numbers like bloom filters or count-min.
    They would have margin of error.




## Problem Statement

Classic Sliding Window to handle time/data structs.
To calculate total days, we need to find the span of the logs and account for overflow.
Formula would be Days = Total Logs (N) / Daily Capcity (10)

