# Here is the tasks for this project

## Jira Daily report

1. Collect information from work items that are updated in the last 24 hours.
2. Group information by epic. If work items are not belong to any epic, put them into "Anything else" section.
3. Sort work items by updated time.
4. Generate report in markdown format. Refer to `.ai/jira-daily-report-template.md` for the format. Whereas:
   1. Heading 1 is the title of the report with format "Daily Report <DD-MMM-YYYY>". `<DD-MMM-YYYY>` is the date format for the report.
   2. Heading 2 is the epic title with format "<Epic KEY> <Epic Status>: <Epic Summary>". If there is not Epic, use "Anything else" as the title.
      1. If the Epic is not updated, and there is any work item in the epic updated, still list the epic in the report.
   3. Heading 3 is the work item title with format "<Task Type> | <Task KEY> <Task Status>: <Task Summary>".
      1. A work item is considered updated if it is created, updated, commented, or work logged in the last 24 hours.
5. Post report to Microsoft Teams channel through Workflow. The Workflow configuration read from environment variables.

The main package source code is located at `cmd/jira-daily-report/main.go`.