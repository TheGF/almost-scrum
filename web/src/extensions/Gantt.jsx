import { FrappeGantt } from 'frappe-gantt-react';
import { React, useContext, useEffect, useState } from 'react';
import Server from '../server';
import UserContext from '../UserContext';
import { Box, Button, ButtonGroup, Text, Stack, HStack, Spacer, Select } from '@chakra-ui/react';
import { useToast } from '@chakra-ui/react';

const viewModes = ['Half Day', 'Day', 'Week', 'Month']

function Gantt(props) {
    const { project } = useContext(UserContext);
    const [tasks, setTasks] = useState([])
    const [viewMode, setViewMode] = useState('Week')
    const [statusFilter, setStatusFilter] = useState([])
    const [boardFilter, setBoardFilter] = useState('')
    const [ownerFilter, setOwnerFilter] = useState('')
    const toast = useToast()

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

    function getOwner(t) {
        const p = t.task.properties
        return p['Owner'] || null
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

    function changeProgress(t, progress) {
        t = findByName(t.name)
        if (t) {
            t.task.properties['Progress'] = `${progress}`
            saveTask(t)
        }
    }

    function chooseBoard(e) {
        setBoardFilter(e && e.target && e.target.value || '')
    }

    function chooseOwner(e) {
        setOwnerFilter(e && e.target && e.target.value || '')
    }

    function toggleStatusFilter(status) {
        if (statusFilter.includes(status)) {
            setStatusFilter(statusFilter.filter(s => s != status))
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
        if (ownerFilter != '' && ownerFilter != p['Owner']) {
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
            dependencies: p['Deps'] && p['Deps'].split(',') || [],
        }
    }

    function detectLoops(ganttTasks) {
        deps = {}
        for (const t of ganttTasks) {
            deps[t.id] = t.dependencies
        }

    }

    function sortByDependency(a, b) {
        const aDependsOnB = a.dependencies.includes(b.id)
        const bDependsOnA = b.dependencies.includes(a.id)

        if (a.dependencies.includes(a.id)) {
            a.dependencies = a.dependencies.filter(d => d.id != a.id)
            toast({
                title: 'Circular Dependency',
                description: `Invalid internal loop in ${a.name}`,
                status: "error",
                isClosable: true,
            })
        }
        if (aDependsOnB && bDependsOnA) {
            b.dependencies = b.dependencies.filter(d => d != a.id)
            toast({
                title: 'Circular Dependency',
                description: `Invalid loop between ${a.name} and ${b.name}`,
                status: "error",
                isClosable: true,
            })
            return 1
        }

        return aDependsOnB ? 1 : -1
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

    const owners = tasks && tasks.length &&
        [...new Set(tasks.map(getOwner).filter(s => s))]


    const statusUI = statusOptions && statusOptions.map(
        s => <Button key={s} isActive={statusFilter.includes(s)}
            onClick={_ => toggleStatusFilter(s)}>
            {s}
        </Button>)

    const boardsUI = boards && boards.map(b => <option key={b} value={b}>
        {b}
    </option>)

    const ownersUI = owners && owners.map(o => <option key={o} value={o}>
        {o}
    </option>)


    const ganttTasks = tasks && tasks.length &&
        tasks.map(getGanttTask)
            .filter(t => t)
            .sort(sortByDependency)
    return ganttTasks ?
        <Stack className="panel1" w="100%" direction="column" spacing={2} p={2}>
            <HStack>
                <ButtonGroup size="sm" >
                    {viewModesUI}
                </ButtonGroup>
                <Spacer />
                <Select maxW="10em" isRequired={false}
                    onChange={chooseBoard} size="sm">
                    <option key="" value={null}></option>
                    {boardsUI}
                </Select>
                <Select maxW="10em" isRequired={false}
                    onChange={chooseOwner} size="sm">
                    <option key="" value={null}></option>
                    {ownersUI}
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
                        onProgressChange={changeProgress}
                    /> : null
                }
            </Box>
        </Stack> :
        <Text>No task in your project has required properties <b>Start</b> and <b>End</b></Text>

}

export default Gantt