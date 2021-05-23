# Boards
Boards group together tasks. Example of boards are the common Scrum *backlog*, *sandbox* and *sprints*.
Boards are listed at the top of the screen close to the __Actions__ menu. Last used boards are visible. Those not recently used are accessible from the dotted button.
A special board is __All__, which is a convenient view to display all tasks from all boards.

## How to create a board
Just select *New Board* in the Actions menu

![](/help/newBoard.png#size=15,align=left)

## The board content
The board contains a filter panel and a list of tasks.

![](/help/board.png#size=90,align=left)
The filter panel includes:
- a button to create a new post
- a search field to filter the tasks
- a button to show/hide the content of tasks
- a button to choose filter criterias

Under the filter panel is the list of tasks, ordered by modification time.

## A Task
A task represents some activity assigned to a person (owner). 
The task owner may usually changes during the project life. 
For instance at creation many tasks are assigned to the product owner and their owner changes during the Sprint planning.

![](/help/taskHeader.png#size=80,align=left)

Each task has an header and a body. The body can hide when the user clicks the collapse button in the filter panel or when the user clicks on the header.

The header shows:
- the numeric id of the task
- the name of the task 
- the tags present in the task (visible only when the body is hidden)
- a button to touch a task and push it to the top of the list
- the percentage progress 
- the board
- the owner 
- a delete button

The user can change the name, the board and the owner.
A user should not reassign a task when he is not the owner: when multiple users are editing the same task, a conflict may occur during Git pull.

![](/help/task.png#size=80,align=left)

The body of the task shows the content of the task in multiple tabs:
- a View tab that shows the Markdown formatted in _HTML_
- a Edit tab to edit the Markdown source
- a Properties tab with the fields that define the tasks, such as the owner, the story points, the status. It is possible to customize the properties (see hacking section of the help)
- a Progress tab where the task may be described in different parts and their completion status
- a Files tab where the task is linked to files stored in the library

Besides the tabs, the body shows the tags present in the task. 
Finally it shows a slider that adjusts the height of the task widget.

### About Editing
Only the owner of a task can edit the task. There is no button Save. 
When the user modifies any content in the task, the update is automatic.
For performance reasons, there is a small delay, so to group together multiple changes.
During the update the __Trash__ button is replaced by a __Save__ button. 

### About Tags
A tag is any token in the task that starts with the character __#__. Tags are specially recognised and highlighted. 
A task property can contain tags (e.g. Status)

### About Indexing
![](/help/searchLorem.png#size=20,align=left)
All content in a task is automatically indexed. 
Then a user can filter the tasks by typing indexed words in the filter panel. The filter is able to suggest words that are relevant.

![](/help/searchStatus.png#size=20,align=left)

When the word starts with _#_, the search is for a tag. 
When the work starts with _@_, the search is for a user

### About Images 
![](/help/imageEdit.png#size=60,align=left)
The Markdown editor supports _copy and paste_ or _drag and drop_ of images. 

![](/help/imageView.png#size=60,align=left)
Once the image has been pasted in the editor, it can be resized and aligned in the HTML view (click on the image).




