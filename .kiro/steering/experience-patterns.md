---
inclusion: auto
---

# 经验模式库

从历次 AI Agent 交互中提炼的可复用模式。

## 指挥经验

- **Hook prompt 祈使句原则**：Hook 的 prompt 必须用明确的祈使句（"编写测试"、"修复错误"），而非询问式（"检查是否需要"）。模糊措辞会让 Agent 倾向于最小行动（只报告不执行）。

## 项目经验

- **收尾任务用 agentStop 而非 postToolUse**：需要在任务完成后统一执行的操作（补测试、格式化、changelog 更新等），优先用 agentStop hook，避免 postToolUse 在每次写操作后重复触发。
