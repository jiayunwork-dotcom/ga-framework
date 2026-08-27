# ga-framework

一个最小可用的遗传算法框架（以 0/1 基因、OneMax 为内置优化目标），演示选择、交叉、变异与精英保留的完整流程。

## 用法

```bash
# 默认参数运行（16 基因、种群 50、100 代）
ga-framework run

# 自定义规模
ga-framework run --genes 32 --size 80 --generations 200 --mutate 0.08 --elite 4 --seed 7
```

输出最优个体的适应度（基因中 1 的个数）与代数。参数非法（规模/基因/代数为 0）时返回受控错误。

## 构建

```bash
go build ./...
go test ./...
```
