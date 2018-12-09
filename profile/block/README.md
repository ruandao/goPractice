
实验了下, ticker 的 C 会在 block profile 中留下足迹, 但是 sleep 操作不好留下足迹

那么使用的时候, 应该选择sleep 还是ticker 呢?

我青睐于使用 sleep 减少,profile 的复杂度, 时间的ticker, 并不是被 block 了,而是时间还没到,还没执行
