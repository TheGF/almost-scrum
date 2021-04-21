import {
    Badge,
    Box, Button, Editable, EditableInput, EditablePreview,
    HStack,
    IconButton, Select, Slider,
    Text, Textarea,
    SliderFilledTrack,
    SliderThumb, SliderTrack, Spacer, Tab, TabList, TabPanel, TabPanels,
    Tabs,
    VStack,
    ButtonGroup
} from "@chakra-ui/react";
import { React, useContext, useEffect, useState } from "react";
import { BsTrash, FiSave, GiRadioactive, MdVerticalAlignTop, GrCompliance } from "react-icons/all";
import T from "../core/T";
import Utils from '../core/utils';
import Server from '../server';
import UserContext from '../UserContext';
import ConfirmChangeOwner from './ConfirmChangeOwner';
import ConfirmDelete from './ConfirmDelete';
import Files from './Files';
import Progress from './Progress';
import Properties from './Properties';
import TaskEditor from './TaskEditor';
import TaskViewer from './TaskViewer';
import MarkdownEditor from '../core/MarkdownEditor';
import { getDefaultToolbarCommands } from 'react-mde'



function Task(props) {
    const { project, info } = useContext(UserContext);
    const { board, name, modTime } = props.info;
    const { boards, users, searchKeys } = props;
    const [compact, setCompact] = useState(props.compact)
    const [saving, setSaving] = useState(false)
    const [task, setTask] = useState(null)
    const [progress, setProgress] = useState('')
    const [openConfirmDelete, setOpenConfirmDelete] = useState(false)
    const [candidateOwner, setCandidateOwner] = useState(null)
    const [height, setHeight] = useState(400)
    const [tabIndex, setTabIndex] = useState(0)


    useEffect(_ => setCompact(props.compact), [props.compact])

    function getTags(task) {
        function extractTags(text) {
            const tags = []
            if (text) {
                const re = /(#\w+)/g
                while (true) {
                    const m = re.exec(text);
                    if (m) { tags.push(m[1]) } else break
                }
            }
            return tags
        }

        let tags = extractTags(task.description)
        for (const value of Object.values(task.properties)) {
            tags = [...tags, ...extractTags(value)]
        }
        for (const part of Object.values(task.parts)) {
            tags = [...tags, ...extractTags(part)]
        }
        return tags
    }

    function updateProgress(task) {
        const progress = task && task.parts && task.parts.length ?
            `${Math.round(100 * task.parts.filter(p => p.done).length / task.parts.length)}%`
            : '-'
        setProgress(progress)
        return task
    }

    function getTask() {
        Server.getTask(project, board, name)
            .then(updateProgress)
            .then(setTask)
    }
    useEffect(getTask, [])

    function touchTask(e) {
        e.stopPropagation()
        Server.touchTask(project, board, name)
            .then(_ => props.onBoardChanged && props.onBoardChanged())
    }

    function saveTask(task) {
        setSaving(true)
        Server.setTaskLater(project, board, name, task)
            .then(_ => setSaving(false))
    }

    function renameTask(title) {
        Server.moveTask(project, board, name, board, title)
            .then(_ => props.onBoardChanged && props.onBoardChanged())
    }

    function deleteTask() {
        Server.deleteTask(project, board, name)
            .then(_ => props.onBoardChanged && props.onBoardChanged())
            .then(_ => setOpenConfirmDelete(false))
    }

    function changeOwner(evt) {
        const newOwner = evt && evt.target && evt.target.value;
        if (newOwner) {
            if (owner == info.loginUser) {
                task.properties['Owner'] = `@${newOwner}`;
                saveTask(task)
                setTask({ ...task })
            } else {
                setCandidateOwner(newOwner);
            }
        }
    }

    function confirmCandidateOwner() {
        task.properties['Owner'] = `@${candidateOwner}`;
        saveTask(task)
        setTask({ ...task })
        setCandidateOwner(null);
    }

    function onBoardChanged(evt) {
        const newBoard = evt && evt.target && evt.target.value;
        newBoard && Server.moveTask(project, board, name, newBoard)
        props.onBoardChanged && props.onBoardChanged()
    }

    if (!task) return null;

    const owner = task && task.properties && task.properties['Owner']
        && task.properties['Owner'].substring(1)
    const readOnly = owner != info.loginUser

    const [id, title] = name && name.split(/\.(.+)/) || ['', 'Something went wrong']
    const userList = users && users.map(u => <option key={u} value={u}>
        {u}
    </option>)
    const boardList = boards && boards.map(b => <option key={b} value={b}>
        {b}
    </option>)

    const mtime = `Last modified: ${Utils.getFriendlyDate(modTime)}`
    const tags = task ? getTags(task).map(tag => <Badge key={tag} colorScheme="purple">
        {tag}
    </Badge>) : null;



    function getHeader() {

        const label = <label>{id}.</label>
        const name = <Editable defaultValue={title} borderWidth="1px" minW="300px"
            borderColor="blue" onSubmit={title => renameTask(title)}>
            <EditablePreview />
            <EditableInput />
        </Editable>
        const compactSwitch = <Spacer onClick={_ => setCompact(!compact)}
            minW="1em" style={{ cursor: 'pointer' }} />
        const tagsGroup = compact ? <HStack h="2em" spacing={2}>{tags}</HStack> : null
        const touchButton = <Button size="sm" title={mtime}><MdVerticalAlignTop onClick={touchTask} /></Button>
        const taskProgress = <span title="Task Progress" style={{ width: '3em', textAlign: 'center' }}>{progress}</span>
        const assignBoard = <Select value={board} title="Assign the Board" w="10em" onChange={onBoardChanged}>
            {boardList}
        </Select>
        const assignOwner = <Select value={owner} title="Assign the Owner" w="10em" onChange={changeOwner}>
            {userList}
        </Select>
        const deleteButton = saving ?
            <IconButton title="Saving..." icon={<FiSave />} /> :
            <IconButton title="Delete the task" icon={<BsTrash />}
                onClick={_ => setOpenConfirmDelete(true)} />

        const confirmChangeOwner = <ConfirmChangeOwner owner={owner} candidateOwner={candidateOwner}
            setCandidateOwner={setCandidateOwner} onConfirm={confirmCandidateOwner} />
        const confirmDelete = <ConfirmDelete isOpen={openConfirmDelete} setIsOpen={setOpenConfirmDelete}
            onConfirm={deleteTask} />



        const header = task.conflictId ? <HStack spacing={3}>
            {label}
            {name}
            {compactSwitch}
            <Badge colorScheme="red"><T>Conflicted Task</T></Badge>
            <GiRadioactive />
            {deleteButton}
            {confirmDelete}
        </HStack> :
            <HStack spacing={3}>
                {label}
                {name}
                {compactSwitch}
                {tagsGroup}
                {touchButton}
                {taskProgress}
                {assignBoard}
                {assignOwner}
                {confirmChangeOwner}
                {deleteButton}
                {confirmDelete}
            </HStack>
        return header
    }


    function handleTabsChange(index) {
        if (index == 0) {
            setTask({ ...task, description: task.description });
        }
        setTabIndex(index)
    }

    const heightSelector = <Slider min={200} max={900} defaultValue="400" w="8em"
        onChangeEnd={setHeight}>
        <SliderTrack>
            <SliderFilledTrack />
        </SliderTrack>
        <SliderThumb />
    </Slider>

    function resolveConflict(target) {
        const lines = task.description.split('\n');
        const description = ''
        let state = 'common'
        for (const line of lines) {
            if (line == '<<<<<<< HEAD') state = 'head'
            else if (line == '=======') state = 'remote'
            else if (line.startsWith('>>>>>>>')) state = 'common'
            else if (state === 'common' ||
                (state === 'head' && target == 'head') ||
                (state === 'remote' && target == 'remote'))
                description += `${line}\n`
        }
        const t = {
            ...task,
            description: description,
            conflictId: '',
        }
        setTask(t)
        saveTask(t)
        props.onBoardChanged && props.onBoardChanged()
    }

    function setDescription(description) {
        const t = {
            ...task,
            description: description,
        }
        setTask(t)
    }

    function markAndSave(opts) {
        const t = {
            ...task,
            description: opts.initialState.text,
            conflictId: '',
        }
        setTask(t)
        saveTask(t)
        props.onBoardChanged && props.onBoardChanged()
    }

    function getBody() {
        if (compact) {
            return ''
        }

        if (task.conflictId) {
            const saveCommand = {
                name: "Save and Resolved",
                icon: () => (
                    <Button>Mark solved and Save <FiSave /></Button>
                ),
                execute: markAndSave
            };

            const toolbarCommands = [...getDefaultToolbarCommands(), ["save-resolve"]]

            return <Box h={height} >
                <Tabs w="100%" index={tabIndex} onChange={handleTabsChange} isLazy>
                    <TabList>
                        <Tab key="conflict"><T>conflict</T></Tab>
                        <Tab key="edit"><T>manual edit - be careful!</T></Tab>
                    </TabList>

                    <TabPanels>
                        <TabPanel key="conflict" padding={0}>
                            <VStack spacing="5" align="left">
                                <Text fontSize="lg">Merge Conflict</Text>
                                <Text colorScheme="red">The task containts a Git conflict.
                                This happens when the same task has been
                                modified by different users at the same time.
                                    <br />
                                    You can solve the problem by using the local copy or the remote one. Or you can edit manually the file.
                                </Text>
                                <ButtonGroup colorScheme="green">
                                    <Button onClick={_ => resolveConflict('head')}>
                                        Use the local copy (Head)
                                    </Button>
                                    <Button onClick={_ => resolveConflict('remote')}>
                                        Use the origin ({task.conflictId})
                                </Button>
                                </ButtonGroup>
                            </VStack>
                        </TabPanel>
                        <TabPanel key="edit" padding={0}>
                            <MarkdownEditor
                                commands={{
                                    "save-resolve": saveCommand
                                }}
                                toolbarCommands={toolbarCommands}
                                value={task.description}
                                height={height + 20}
                                onChange={setDescription}
                                disablePreview={true}
                                paste={null}
                            />
                        </TabPanel>
                    </TabPanels>
                </Tabs>
            </Box>
        }


        return <Box h={height} ><HStack spacing={3} >
            <Tabs w="100%" index={tabIndex} onChange={handleTabsChange} isLazy>
                <TabList>
                    <Tab key="edit"><T>{readOnly ? 'view' : 'edit'}</T></Tab>
                    <Tab key="progress"><T>progress</T></Tab>
                    <Tab key="files"><T>files</T></Tab>
                    <Spacer key="spacer" />
                    <HStack h="2em" spacing={2} key="tags">{tags}</HStack>
                    <div width="2em" />
                    <Spacer maxWidth="2em" />
                    {heightSelector}
                </TabList>

                <TabPanels>
                    <TabPanel key="edit" padding={0}>
                        <TaskEditor task={task} saveTask={saveTask} users={users} height={height}
                            readOnly={readOnly} />
                    </TabPanel>
                    <TabPanel key="progress" >
                        <Progress task={task} readOnly={readOnly} height={height} 
                            saveTask={task => {
                                saveTask(task);
                                updateProgress(task);
                            }} />
                    </TabPanel>
                    <TabPanel>
                        <Files task={task} saveTask={saveTask} readOnly={readOnly}
                            height={height} />
                    </TabPanel>
                </TabPanels>
            </Tabs>
        </HStack> </Box>
    }


    if (!task) {
        return null
    }
    return <Box p={1} w="100%" borderWidth="3px" overflow="hidden">
        {getHeader()}
        {getBody()}
    </Box>
}

export default Task;