Go 遗传算法命令行框架，内置 0/1 基因 OneMax：ga-framework run 按种群规模、变异率和精英保留迭代并打印最优适应度；serve 提供 /api/run 做同一次求解。

# ga-framework

Genetic algorithm optimization framework with configurable selection strategies,
crossover/mutation operators, island migration, niching, and convergence
termination. Exposes HTTP API for parameter sweep and single-run execution.

## Build / Run / Test

```bash
go build -o ga-framework .
./ga-framework serve --addr :8080
./ga-framework run --size 100 --genes 32 --generations 200
go test ./...
```

## Evaluation Image

Evaluation-specific files (do not overwrite project Dockerfile/README):

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md` (this file)

Build and verify in container:

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
