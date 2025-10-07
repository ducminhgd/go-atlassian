# Here is the tasks for this project

## Jira Daily report

1. Jira has 3 levels of work items: Epic, Task/Story/Bug, Sub-task. A task can belong to an epic, and a sub-task can belong to a task.
2. Collect information from work items that are updated in the last 24 hours.
3. Group information by epic. If work items are not belong to any epic, put them into "Anything else" section.
4. Sort work items by updated time.
5. Generate report in markdown format or AdaptiveCard format
   1. . Refer to `.ai/jira-daily-report-template.md` for the Markdown format. Whereas:
      1. Heading 1 is the title of the report with format "Daily Report <DD-MMM-YYYY>". `<DD-MMM-YYYY>` is the date format for the report.
      2. Heading 2 is the level 2 work items (Epic) title with format "<Epic KEY> | <Epic Status> | <Epic Summary>". Use "Anything else" as the title for orphan work items.
         1. If the Epic is not updated, and there is any work item in the epic updated, still list the epic in the report.
      3. Heading 3 are the level 2 work items, title with format "<Task Type> | <Task KEY> | <Task Status> | <Task Summary>".
         1. A work item is considered updated if it is created, updated, commented, or work logged in the last 24 hours.
      4. Heading 4 are the level 3 work items (sub-tasks) belong to the work item in Heading 3.
      5. If a level-2 or level-3 work item is updated but not their parent, go fetch their parent and treat them as updated.
   2. If AdaptiveCard format is chosen, convert the markdown to AdaptiveCard format. https://github.com/OfficeDev/Microsoft-Teams-Samples/tree/main/tools/message-card-to-ac-transformation
      1. Heading 1 (title of the report) is a TextBlock with extralarge size.
      2. Heading 2 (level 1 work items, Epic or Anything else) is a TextBlock with large size and bold weight.
      3. Heading 3 (level 2 work items, Task/Story/Bug) is a TextBlock with medium size and bold weight.
      4. Heading 4 (level 3 work items, Sub-task) is a TextBlock with small size and bold weight, underlined.
      5. Comments or worklogs are TextBlocks with contents are quoted, number bulletings
6. Post report to Microsoft Teams channel through Workflow. The Workflow configuration read from environment variables.

The main package source code is located at `cmd/jira-daily-report/main.go`.