/*
广度优先算法走迷宫
1.用循环创建二维slice
2.使用slice来实现队列
3.使用fscanf读取文件
4.对point进行抽象
*/
package main

import (
	"fmt"
	"os"
)

func main() {
	maze := readMaze("./src/maze/maze.txt")
	//showMaze(maze)
	steps, minStep, minSteps := walk(maze, point{0, 0}, point{len(maze) - 1, len(maze[0]) - 1})
	showMaze(steps)
	fmt.Println(minStep)
	for i, p := range minSteps {
		fmt.Printf("%d -> %v \n", i, p)
	}
}

// read file of maze.txt
func readMaze(fileName string) [][]int {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// defend row and col
	var row, col int
	// read data from file ,因为读进来，会给row，col赋值，所以要用地址
	fmt.Fscanf(file, "%d %d", &row, &col)
	// create an slice,定义一个大 slice ,大slice每一个element，都是一个小slice
	maze := make([][]int, row)
	for i := range maze {
		maze[i] = make([]int, col)
		for j := range maze[i] {
			fmt.Fscanf(file, "%d", &maze[i][j])
		}
	}

	return maze
}

// defined point struct
type point struct {
	i, j int
}

// 为point结构定义点和点的加法
// defined function addition for point struct
func (p point) added(r point) point {
	return point{p.i + r.i, p.j + r.j}
}

// 获取当前点在maze中的值，同时还不能越界
// 第二个返回值，常用于map中是否有值，channel是否被close了
func (p point) at(grid [][]int) (int, bool) {
	// 判断是否越界 往上，往下是否走出了row 行数 越界
	if p.i < 0 || p.i >= len(grid) {
		return 0, false
	}
	//判断是否越界 往左，往右是否走出了col 列数 越界
	if p.j < 0 || p.j >= len(grid[p.i]) {
		return 0, false
	}
	return grid[p.i][p.j], true
}

/*
        列
0 1 0 0 0 行
0 0 0 1 0
0 1 0 1 0
1 1 1 0 0
0 1 0 0 1
0 1 0 0 0
定义一下四个方向
*/

var dirs = [4]point{
	{-1, 0}, //上 行减一，列不动
	{0, -1}, //左 行不动，列减一
	{1, 0},  //下 行加一，列不动
	{0, 1},  //右 行不动，列加一
}

// 走迷宫
// maze 地图
// start 地图起始坐标
// end 地图结束坐标
func walk(maze [][]int, start, end point) ([][]int, int, map[int]point) {
	//维护从start，到end，一共走了多少格才走到终点。也就是我们探索过的点，全部放进去这个slice
	steps := make([][]int, len(maze))
	for i := range steps {
		steps[i] = make([]int, len(maze[i]))
	}
	//定义point队列,用于保存第二个循环找到的探索的点，在里面排队去探索
	Q := []point{start}
	//当队列不为空的时候我们去探索
	for len(Q) > 0 {
		//开始探索 拿出当前的点
		currentElement := Q[0]
		Q = Q[1:]
		// 以currentElement为节点，去进行上左下右四个方向的探索
		for _, dir := range dirs {
			//开始每个方向探索节点,返回探索到的节点
			nextElement := currentElement.added(dir)
			// maze at next is 0 迷宫的下一个探索点是0才能走
			// and steps at next is 0 如果steps中有这个点，说明已经走过了
			// and next!=start 如果等于start说明走回起点去了也不行

			//获取当前探索的坐标点，在maze里的值
			value, ok := nextElement.at(maze)
			// 如果越界和value=1 撞墙 就进行下一个方向探索
			if !ok || value == 1 {
				continue
			}
			// 检查当前坐标点在steps有返回值，说明已经走过了
			value, ok = nextElement.at(steps)
			if !ok || value != 0 {
				continue
			}
			//如果当前坐标点等于start说明又回到了原点
			if nextElement == start {
				continue
			}
			// 去steps slice中设置探索坐标点的值，等于当前坐标点的值+1
			value, _ = currentElement.at(steps)
			steps[nextElement.i][nextElement.j] = value + 1
			// 然后将这个要探索的坐标点加入队列
			Q = append(Q, nextElement)

			//如果到终点了就退出所有循环
			if nextElement == end {
				goto endFor
			}
		}
	}
endFor:
	// 获取最短 步数
	minStep, _ := end.at(steps)
	// 定义最小路径map
	minSteps := map[int]point{
		minStep: end,
	}
	// 反着探索一遍只要nextEle=curEle-1就放进map
	Rq := []point{end}
	// 当步数和走过的点数相等的时候退出，不严谨的话直接不加条件
	for len(minSteps) < (minStep + 1) {
		// 获取当前节点
		currentElement := Rq[0]
		// 将当前节点移除队列
		Rq = Rq[1:]
		// 循环四个方向
		for _, dir := range dirs {
			//获取要探索的节点
			nextElement := currentElement.added(dir)
			//获取要探索节点的值
			nextVal, ok := nextElement.at(steps)
			// 如果越界和value=0 撞墙 就进行下一个方向探索
			if !ok || ((nextVal == 0) && (nextElement != start)) {
				continue
			}
			//获取当前point的值
			currentVal, ok := currentElement.at(steps)
			// 如果当前point的值不等于 探索到的节点值+1 进行下一个方向
			if currentVal -= 1; currentVal != nextVal {
				continue
			}
			// 下一次探索以这个point为中心点
			Rq = append(Rq, nextElement)
			//  把探索到的点，添加到map中
			minSteps[nextVal] = nextElement
			//如果到终点了就退出所有循环
			if nextElement == start {
				goto endStep
			}
		}
	}
endStep:
	return steps, minStep, minSteps
}

// 显示迷宫
func showMaze(maze [][]int) {
	for _, row := range maze {
		for _, col := range row {
			fmt.Printf("%3d ", col)
		}
		fmt.Println()
	}
}
