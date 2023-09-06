/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Test struct {
	A string
	B string
}

func main() {
	h := server.Default()
	h.StaticFS("/", &app.FS{Root: "./", GenerateIndexPages: true})

	count := 0

	h.POST("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.Response.Header.Set("X-Tt-Logid", strconv.Itoa(count))
		count++
		ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
	})

	h.GET("/json", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, &Test{
			A: "aaa",
			B: "bbb",
		})
	})

	h.GET("/redirect", func(c context.Context, ctx *app.RequestContext) {
		ctx.Redirect(consts.StatusMovedPermanently, []byte("http://www.google.com/"))
	})

	v1 := h.Group("/v1")
	{
		v1.GET("/hello/:name", func(c context.Context, ctx *app.RequestContext) {
			fmt.Fprintf(ctx, "Hi %s, this is the response from Hertz.\n", ctx.Param("name"))
		})
	}

	go func() {
		cli, _ := client.NewClient(client.WithResponseBodyStream(true))
		time.Sleep(2 * time.Second)
		for i := 0; i < 500; i++ {
			req := protocol.AcquireRequest()
			resp := protocol.AcquireResponse()
			req.SetMethod(consts.MethodPost)
			req.SetRequestURI("http://127.0.0.1:8888/ping")
			err := cli.Do(context.Background(), req, resp)
			if err != nil {
				fmt.Println("error:" + err.Error())
				return
			}
			fmt.Println(resp.StatusCode())
			resp.Reset()
			time.Sleep(200 * time.Millisecond)
		}

	}()

	h.Spin()
}
