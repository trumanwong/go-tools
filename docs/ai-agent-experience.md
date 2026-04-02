# AI Agent 交互经验库

记录与 AI Agent 协作过程中的决策案例和经验模式。

### #1 Hook prompt 措辞决定 Agent 行为边界

- **时间**：2026-04-01
- **背景**：用户要求给 crawler 的 Request 结构体增加 PostForm 字段并修改 Send 方法，Agent 完成代码修改后，post-generation-review hook 触发了测试用例检查步骤，但 Agent 只做了评估没有实际编写测试
- **问题**：用户期望 Agent 在修改代码后自动补写测试用例，但 hook 的 prompt 措辞是"检查是否需要更新"而非"编写或更新测试"，导致 Agent 只做检查不做执行
- **解法**：讨论了两种 hook 方案：1) postToolUse + write 类型，每次写文件后触发（粒度细但可能重复触发）；2) agentStop 类型，任务结束后统一触发一次（更干净）。推荐 agentStop 方案
- **结论**：Hook 的 prompt 措辞必须用祈使句明确指令（"编写测试"而非"检查是否需要"），否则 Agent 会倾向于保守执行
- **分类**：指挥经验 / 项目经验
- **思维模式**：
  - Hook prompt 是给 Agent 的指令，措辞要用"做什么"而非"看看是否需要做什么"，模糊措辞会导致 Agent 选择最小行动
  - 对于自动化流程，agentStop 比 postToolUse 更适合做"收尾检查+补充"类任务，避免重复触发
