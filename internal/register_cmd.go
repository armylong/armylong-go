package internal // 包名=目录名，适配internal/根目录

import (
	"os"

	"github.com/armylong/armylong-go/internal/cmd" // 导入handler所在包
	"github.com/spf13/cobra"
)

// RootCmd 根命令（Cobra标准写法，Use定义命令名）
var RootCmd = &cobra.Command{
	Use:   "trae-llm-testing-dogfood-2-553",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// 无参数时打印帮助（提升体验）
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute 执行根命令（暴露给main.go调用）
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// 根命令本地标志
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// 初始化时注册所有子命令
	RegisterCmd()
}

// RegisterCmd 集中注册所有子命令（修正所有Cobra语法错误）
func RegisterCmd() {
	// ==================================================
	// 步骤1：创建命令实例（Use替代Name，删除无效的Flags数组）
	demoCmd := &cobra.Command{
		Use:   "demo [username]", // Cobra核心：Use定义命令名+使用方式
		Short: "演示参数接收",
		Args:  cobra.MaximumNArgs(1),
		Run:   cmd.DemoHandler, // 调用internal/cmd下的执行逻辑
	}
	// 步骤2：为命令添加标志（Cobra标准方式，替代Flags数组）
	demoCmd.Flags().BoolP("enable", "e", false, "是否启用功能（布尔类型）")
	demoCmd.Flags().StringP("message", "m", "", "自定义消息（字符串类型，必填）")
	_ = demoCmd.MarkFlagRequired("message") // 标记必填（替代Required:true）
	demoCmd.Flags().IntP("age", "a", 18, "用户年龄（整数类型）")
	demoCmd.Flags().StringSliceP("hobby", "H", []string{}, "爱好（可多次指定，如 -H 篮球 -H 编程）")
	// 步骤3：将命令挂载到根命令
	RootCmd.AddCommand(demoCmd)

	// ==================================================
	uvRecalcCmd := &cobra.Command{
		Use:   "UvLampStatisticsRecalculate", // 替换Name→Use
		Short: "紫外线灯统计重算",
		Run:   cmd.UvLampStatisticsRecalculateHandler,
	}
	// 添加标志
	uvRecalcCmd.Flags().String("date", "", "指定日期，格式：2006-01-02（可选，默认今天，与start-date/end-date互斥）")
	uvRecalcCmd.Flags().String("start-date", "", "开始日期，格式：2006-01-02（可选，与date互斥）")
	uvRecalcCmd.Flags().String("end-date", "", "结束日期，格式：2006-01-02（可选，与start-date配合使用，不指定则等于start-date）")
	uvRecalcCmd.Flags().Int64("tenant-id", 0, "门店ID（可选，不指定则扫描全部门店）")
	RootCmd.AddCommand(uvRecalcCmd)

	// ==================================================
	redisGetDataCmd := &cobra.Command{
		Use:   "GetRedisData [cache_key]", // 替换Name→Use
		Short: "从Redis缓存中读取数据",
		Run:   cmd.GetRedisData,
	}
	RootCmd.AddCommand(redisGetDataCmd)

	// ==================================================
	todoCmd := &cobra.Command{
		Use:   "todo [task_type]", // 替换Name→Use
		Short: "任务管理",
		Long: `任务管理CLI工具，支持以下子命令:
  create   创建任务 (参数: --title, --desc, --sort, --expire_at)
  get      获取任务详情 (参数: --task_id)
  sort     更新任务排序 (参数: --task_id, --sort)
  complete 标记任务完成 (参数: --task_id)
  expire   检测并标记过期任务 (无需参数)`,
		Run: cmd.TodoHandler,
	}
	todoCmd.Flags().Int64P("task_id", "", 0, "任务ID（可选）")
	todoCmd.Flags().StringP("title", "", "", "任务标题（create时必填）")
	todoCmd.Flags().StringP("desc", "", "", "任务描述（create时必填）")
	todoCmd.Flags().Int64P("sort", "", 0, "任务排序值，数字越大越靠前（可选）")
	todoCmd.Flags().StringP("expire_at", "", "", "过期时间，格式：2006-01-02 15:04:05（可选）")
	RootCmd.AddCommand(todoCmd)

}
