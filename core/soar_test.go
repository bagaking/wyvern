package core_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bagaking/wyvern/core"
	"github.com/bagaking/wyvern/core/flaps"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
	"time"
)

// config - 测试配置, 包含 print 插件
var config = `
soars:
  - name: soar1
    flaps:
      - name: flap1
        plugin: print
        pluginConfig: 
          msg: "hello world"
`

// captureStdout - 捕获 stdout
func captureStdout(f func()) string {
	// 创建一个管道
	r, w, _ := os.Pipe()
	// 保存原来的 stdout
	stdout := os.Stdout
	// 将 stdout 指向管道的写入端
	os.Stdout = w
	// 执行函数
	f()
	// 将 stdout 指向原来的 stdout
	os.Stdout = stdout
	// 关闭管道的写入端
	w.Close()
	// 从管道的读取端读取数据
	var buf bytes.Buffer
	io.Copy(&buf, r)
	// 返回数据
	return buf.String()
}

// 用 json 序列化并打印
func printJSON(v interface{}) {
	b, e := json.Marshal(v)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(string(b))
	}
}

// TestWyvern_LoadFromConfig - 测试 Wyvern 的 LoadFromConfig 方法
func TestWyvern_LoadFromConfig(t *testing.T) {
	// 从配置文件中加载所有的 Soar
	w := core.NewWyvern()
	// 加载配置
	cfg, err := core.LoadWyvernConfig(config)

	// 断言没有错误, 失败则终止测试
	assert.NoError(t, err)

	// 加载 Soar1
	id, err := w.LoadFromConfig(cfg, "soar1")
	printJSON(id)

	// 断言没有错误, 失败则终止测试
	assert.NoError(t, err)

	// 断言 Soar 清单不为空, 失败则终止测试
	assert.NotEmpty(t, w.Soars)

	// 打印 json 序列化的 Soar 清单
	printJSON(w.Soars)
	flap1 := w.Soars[id].FindFlapByName("flap1")
	// 打印 json 序列化的 flap1
	printJSON(flap1)
	// 断言 Soar 清单中有 flap1
	assert.NotEmpty(t, flap1)
	// 断言 flap1 的插件是 print
	assert.Equal(t, "print", w.Soars[id].FindFlapByName("flap1").Action.Plugin())
	// 获取 flap1 的插件配置
	fp := w.Soars[id].FindFlapByName("flap1").Action.(*flaps.FlapPrint)
	// 打印 json 序列化的 flap1 的插件配置
	printJSON(fp)

	// 断言 flap1 的插件配置是 {msg: "hello world"}
	assert.Equal(t, "hello world", fp.Msg)
}

// TestWyvern_Run - 测试 Wyvern 的 Run 方法
func TestWyvern_Run(t *testing.T) {
	// 从配置文件中加载所有的 Soar
	w := core.NewWyvern()
	// 加载配置
	cfg, err := core.LoadWyvernConfig(config)

	// 断言没有错误, 失败则终止测试
	assert.NoError(t, err)

	// 加载 Soar1
	id, err := w.LoadFromConfig(cfg, "soar1")
	printJSON(id)

	// 断言没有错误, 失败则终止测试
	assert.NoError(t, err)

	// 断言 Soar 清单不为空, 失败则终止测试
	assert.NotEmpty(t, w.Soars)

	// 开始捕获 stdout
	p := captureStdout(func() {
		// 运行 soar1
		err = w.Run(context.Background(), id)
		// 等待一段时间
		time.Sleep(time.Second)

		// 断言没有错误
		assert.NoError(t, err)
	})

	// 断言输出的内容
	assert.Equal(t, "hello world", p)
}
