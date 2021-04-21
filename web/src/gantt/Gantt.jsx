import { FrappeGantt } from 'frappe-gantt-react';
import { React, useContext, useEffect, useState } from 'react';
import Server from '../server';
import UserContext from '../UserContext';

function Gantt(props) {
    const { project } = useContext(UserContext);
    const [tasks, setTasks] = useState([])

    function getTasks() {
        Server.getGanttTasks(project)
            .then(tasks => setTasks(tasks || []))
    }
    useEffect(getTasks, [])

    const ganttTasks = tasks.map(t => ({
        id: t.name,
        name: t.name,
        start: t.task.properties['Start'] || new Date(),
        end: t.task.properties['End'] || new Date(),
    }))

    let d1 = new Date();
    let d2 = new Date();
    d2.setDate(d2.getDate() + 5);
    let d3 = new Date();
    d3.setDate(d3.getDate() + 8);
    let d4 = new Date();
    d4.setDate(d4.getDate() + 20);
    const tasks3 = [
        {
            id: "Task 1",
            name: "Task 1",
            start: d1,
            end: d2,
            progress: 10,
            dependencies: ""
        },
        {
            id: "Task 2",
            name: "Task 2",
            start: d3,
            end: d4,
            progress: null
            // dependencies: "Task 1"
        },
        {
            id: "Task 3",
            name: "Redesign website",
            start: new Date(),
            end: d4,
            progress: 0
            // dependencies: "Task 2, Task 1"
        }
    ];


    return <div>
        <FrappeGantt
            tasks={tasks3}
        />
    </div>

}

export default Gantt