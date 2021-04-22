import { FrappeGantt } from 'frappe-gantt-react';
import { React, useContext, useEffect, useState } from 'react';
import Server from '../server';
import UserContext from '../UserContext';
import { Box, Button, ButtonGroup, Text, Stack, HStack, Spacer, Select } from '@chakra-ui/react';

const viewModes = ['Half Day', 'Day', 'Week', 'Month']

function Gantt(props) {
    const { project } = useContext(UserContext);
    const [tasks, setTasks] = useState([])
    const [viewMode, setViewMode] = useState('Week')
    const [statusFilter, setStatusFilter] = useState([])
    const [boardFilter, setBoardFilter] = useState('')

    function ganttRef(r) {
        if (r) {
            const h = r._svg.current.getAttribute('height')
            r._svg.current.setAttribute('height', h - 120)
        }
    }

    function getTasks() {
        Server.getGanttTasks(project)
            .then(tasks => setTasks(tasks || []))
    }
    useEffect(getTasks, [])

    function getStatusOption(t) {
        const p = t.task.properties
        return p['Status'] || null
    }

    function findByName(name) {
        return tasks.filter(t => t.name == name)[0]
    }

    function saveTask(t) {
        Server.setTask(project, t.board, t.name, t.task)
    }

    function changeDates(t, start, end) {
        t = findByName(t.name)
        if (t) {
            t.task.properties['Start'] = start.toISOString()
            t.task.properties['End'] = end.toISOString()
            saveTask(t)
        }
    }

    function chooseBoard(e) {
        setBoardFilter(e && e.target && e.target.value || '')
    }

    function toggleStatusFilter(status) {
        if (statusFilter.includes(status)) {
            setStatusFilter(statusFilter.filter(s=>s != status))
        } else {
            setStatusFilter([...statusFilter, status])
        }
    }

    function getGanttTask(t) {
        const p = t.task.properties
        const day = 60 * 60 * 24 * 1000;
        const status = p['Status']

        if (statusFilter.length && !statusFilter.includes(status)) {
            return null
        }
        if (boardFilter != '' && boardFilter != t.board) {
            return null
        }

        let start = new Date(p['Start'])
        let end = new Date(p['End'])
        start = isNaN(start.getTime()) ? new Date() : start
        end = isNaN(end.getTime()) ? new Date(start.getTime() + 2 * day) : end
        const progress = p['Progress'] || '10'

        return {
            id: t.name,
            name: t.name,
            start: start,
            end: end,
            progress: progress,
        }
    }

    const viewModesUI = viewModes.map(
        s => <Button key={s} isActive={s == viewMode}
            onClick={_ => setViewMode(s)}>
            {s}
        </Button>)

    const statusOptions = tasks && tasks.length &&
        [...new Set(tasks.map(getStatusOption).filter(s => s))]

    const boards = tasks && tasks.length &&
        [...new Set(tasks.map(t => t.board))]

    const statusUI = statusOptions && statusOptions.map(
        s => <Button key={s} isActive={statusFilter.includes(s)}
            onClick={_ => toggleStatusFilter(s)}>
            {s}
        </Button>)

    const boardsUI = boards && boards.map(b => <option key={b} value={b}>
        {b}
    </option>)


    const ganttTasks = tasks && tasks.length && tasks.map(getGanttTask).filter(t => t)
    return ganttTasks ?
        <Stack className="panel1" w="100%" direction="column" spacing={2} p={2}>
            <HStack>
                <ButtonGroup size="sm" >
                    {viewModesUI}
                </ButtonGroup>
                <Spacer />
                <Select maxW="12em" isRequired={false} onChange={chooseBoard}>
                    <option key="" value={null}></option>
                    {boardsUI}
                </Select>
                <ButtonGroup size="sm" >
                    {statusUI}
                </ButtonGroup>
            </HStack>
            <Box w="100%" borderWidth={1}>
                {
                    ganttTasks.length ? <FrappeGantt
                        viewMode={viewMode}
                        tasks={ganttTasks}
                        ref={ganttRef}
                        // onClick={task => console.log(task)}
                        onDateChange={changeDates}
                        onProgressChange={(task, progress) => console.log(task, progress)}
                        onTasksChange={tasks => console.log(tasks)}
                    /> : null
                }
            </Box>
        </Stack> :
        <Text>No task in your project has required properties <b>Start</b> and <b>End</b></Text>

}

export default Gantt