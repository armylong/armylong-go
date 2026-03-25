package internal // 包名=目录名，适配internal/根目录

import (
	"github.com/armylong/armylong-go/internal/cmd" // 导入handler所在包
	"github.com/armylong/go-library/service/command"
	"github.com/spf13/cobra"
)

// RegisterCmd 集中注册所有子命令（修正所有Cobra语法错误）
func RegisterCmd(command command.BaseCommand) {
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
	command.AddCliCommand(todoCmd)

	command.AddCliCommand(&cobra.Command{
		Use:   "todo2 [task_type] task_id title desc sort expire_at",
		Short: "任务管理",
		Run:   cmd.TodoHandler,
	})

}
