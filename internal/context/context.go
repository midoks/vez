package context

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/flamego/flamego"
	"github.com/flamego/template"

	"github.com/midoks/vez/internal/mgdb"
)

var (
	mg     []mgdb.ContentBid
	mux    sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
)

func InitMG() {
	ctx, cancel = context.WithCancel(context.Background())

	// 启动后台goroutine定期更新数据
	go func() {
		ticker := time.NewTicker(5 * time.Minute) // 5分钟更新一次
		defer ticker.Stop()

		// 初始化数据
		updateRandomContent()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				updateRandomContent()
			}
		}
	}()
}

// StopMG 停止后台goroutine
func StopMG() {
	if cancel != nil {
		cancel()
	}
}

// updateRandomContent 更新随机内容
func updateRandomContent() {
	newMg, err := mgdb.ContentRandNum(10)
	if err != nil {
		fmt.Printf("Error updating random content: %v\n", err)
		return
	}

	mux.Lock()
	mg = newMg
	mux.Unlock()
}

// getRandomContent 安全地获取随机内容
func getRandomContent() []mgdb.ContentBid {
	mux.RLock()
	defer mux.RUnlock()

	if len(mg) == 0 {
		// 如果没有数据，尝试立即获取
		mux.RUnlock()
		updateRandomContent()
		mux.RLock()
	}

	// 返回副本以避免并发问题
	result := make([]mgdb.ContentBid, len(mg))
	copy(result, mg)
	return result
}

func Contexter() flamego.Handler {
	return func(c flamego.Context, t template.Template, d template.Data) {
		d["Newsest"] = getRandomContent()
	}
}
