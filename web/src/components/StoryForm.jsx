import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Badge from 'react-bootstrap/Badge';
import Button from 'react-bootstrap/Button';
import ButtonGroup from 'react-bootstrap/ButtonGroup';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';
import Card from 'react-bootstrap/Card';
import Form from 'react-bootstrap/Form';
import ListGroup from 'react-bootstrap/ListGroup';
import Modal from 'react-bootstrap/Modal';
import Tab from 'react-bootstrap/Tab';
import Tabs from 'react-bootstrap/Tabs';
import ToggleButton from 'react-bootstrap/ToggleButton';
import { useForm } from 'react-hook-form';
import { BiArchive } from 'react-icons/bi';
import { CgMoveTask } from 'react-icons/cg';
import { GiSave } from 'react-icons/gi';
import { getConfig, loginWhenUnauthorized } from './axiosUtils';
import { pointsUnitOptions, statusOptions } from './consts';
import Server from './Server';

function MainForm(props) {
    const { form, project, store, story, moveTo, reload } = props;
    const { register, watch, setValue } = form;

    const [pointsUnit, setPointsUnit] = useState(localStorage.pointsUnit || 'fibo')
    const [users, setUsers] = useState([])
    const [showConfirmArchive, setShowConfirmArchive] = useState(false);

    function ConfirmArchive(props) {
        const {show, setShow} = props;
        return <Modal show={show}>
            <Modal.Header>
                <Modal.Title>Confirmation</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                Do you want to archive this story?
            </Modal.Body>
            <Modal.Footer>
                <Button variant="primary" onClick={
                    _ => Server.moveStory(project, store, story, 'archive')
                        .then(reload)}>
                    Confirm
                </Button>
                <Button variant="secondary" onClick={
                    _=>setShow(false)}>
                    Cancel
                </Button>
            </Modal.Footer>
        </Modal>;
    }
    

    function fetchUsers() {
        axios.get(`/api/v1/projects/${project}/users`, getConfig())
            .then(r => setUsers(r.data))
            .catch(loginWhenUnauthorized)
    }
    useEffect(fetchUsers, [])

    function updatePointsUnit(pointsUnit) {
        localStorage.pointsUnit = pointsUnit;
        setPointsUnit(pointsUnit);
    }

    return <Card body>
        <Form>
            <ConfirmArchive show={showConfirmArchive} setShow={setShowConfirmArchive}/>
            <ButtonToolbar className="float-right">
                {moveTo &&
                    <Button
                        className="mr-1 mb-2"
                        onClick={_ => Server.moveStory(project, store, story, moveTo).then(reload)}>
                        <CgMoveTask /> Move to {moveTo}
                    </Button>}
                <Button className="mr-1 mb-2"
                        onClick={_ => setShowConfirmArchive(true)}>
                    <BiArchive /> Archive
                </Button>
            </ButtonToolbar>
            <Form.Group controlId="Description">
                <Form.Label>Description</Form.Label>
                <Form.Control name="description" as="textarea" rows={3} placeholder="Description"
                    ref={register} />
            </Form.Group>

            <Form.Group controlId="Points">
                <Form.Label>Points</Form.Label>&nbsp;
            <a href="#" onClick={() => updatePointsUnit('fibo')}>fibo</a>&nbsp;
            <a href="#" onClick={() => updatePointsUnit('linear')}>linear</a>
                <br />
                <ButtonGroup toggle name="points" ref={register}>
                    {pointsUnitOptions[pointsUnit].map((number, idx) => (
                        <ToggleButton
                            key={idx}
                            type="radio"
                            variant="secondary"
                            checked={number == watch('points')}
                            size="sm"
                            value={number}
                            onChange={e => setValue('points', parseInt(e.target.value), { shouldValidate: true })}
                        >
                            {number}
                        </ToggleButton>

                    ))}
                </ButtonGroup>
            </Form.Group>

            <Form.Group controlId="Status">
                <Form.Label>Status</Form.Label>
                <Form.Control as="select" name="status" ref={register}>
                    {statusOptions.map(op =>
                        <option key={op}>{op}</option>)}
                </Form.Control>
            </Form.Group>

            <Form.Group controlId="Owner">
                <Form.Label>Owner</Form.Label>
                <Form.Control as="select" name="owner" ref={register}>
                    {users.map(op =>
                        <option key={op}>{op}</option>)}
                </Form.Control>
            </Form.Group>
        </Form>
    </Card>;

}

function TasksForm(props) {
    function EditTaskModal(props) {
        const { form } = props;
        function updateTask(data) {
            const tasks = form.watch('tasks') || [];
            if (editTaskId == 'New')
                form.setValue('tasks', [...tasks, data.taskText], { shouldValidate: true });
            else {
                tasks[editTaskId] = data.taskText;
                form.setValue('tasks', tasks, { shouldValidate: true });
            }

            setEditTaskId(false);
        }

        const text = Number.isFinite(editTaskId) ? watch('tasks')[editTaskId] : '';
        const { register, handleSubmit } = useForm({ defaultValues: { taskText: text } });

        return <Modal show={editTaskId !== false}>
            <Modal.Header>
                <Modal.Title>Edit Task</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <Form>
                    <Form.Group controlId="formBasicEmail">
                        <Form.Label>Task</Form.Label>
                        <Form.Control as="textarea" rows={3} placeholder="Enter the task"
                            name="taskText" ref={register} />
                    </Form.Group>
                </Form>
            </Modal.Body>
            <Modal.Footer>
                <Button variant="primary" onClick={handleSubmit(updateTask)}>
                    Save
                </Button>
                <Button variant="secondary" onClick={_ => setEditTaskId(false)}>
                    Cancel
                </Button>
            </Modal.Footer>
        </Modal>;
    }
    const { form } = props;
    const { watch, setValue } = form;
    const [editTaskId, setEditTaskId] = useState(false)


    const tasksList = (watch('tasks') || []).map((t, idx) => {
        function setTaskState(idx, done) {
            const tasks = watch('tasks');
            tasks[idx] = tasks[idx].replace(/!$/, '') + (done ? '!' : '');
            setValue('tasks', tasks, { shouldValidate: true });
        }

        const done = t.endsWith('!')
        return <ListGroup.Item key={idx} onClick={_ => setEditTaskId(idx)} style={{ cursor: 'pointer' }}>
            {t}
            <ButtonGroup toggle className="float-right">
                <ToggleButton type="radio" variant="primary" active={!done} onClick={e => {
                    e.stopPropagation();
                    if (done) {
                        setTaskState(idx, false);
                    }
                }} size="sm" checked={!done}>TODO</ToggleButton>
                <ToggleButton type="radio" variant="primary" active={done} onClick={e => {
                    e.stopPropagation();
                    if (!done) {
                        setTaskState(idx, true);
                    }

                }} size="sm" checked={done}>DONE</ToggleButton>
            </ButtonGroup>
        </ListGroup.Item>;
    });

    return <Card body>
        <EditTaskModal form={form} />
        <Form>
            <br />
            <Button onClick={_ => setEditTaskId('New')}>New Task</Button>
            <ListGroup>
                {tasksList}
            </ListGroup>
        </Form>
    </Card>;
}

function StoryForm(props) {
    const { project, store, story, form, pendingWrite } = props;
    const { register } = form;

    useEffect(() => register('points'), [register]);
    useEffect(() => register('tasks'), [register]);

    return <>
        <span className="float-right">
            {pendingWrite ? <GiSave /> : null}
            <Badge>{project}/{store}/{story}</Badge>
        </span>
        <Tabs defaultActiveKey="main" >
            <Tab title="Main" eventKey="main">
                <MainForm {...props} />
            </Tab>
            <Tab title="Tasks" eventKey="tasks">
                <TasksForm {...props} />
            </Tab>
            <Tab title="Time" eventKey="time" >

            </Tab>
            <Tab title="Library" eventKey="library" >

            </Tab>
        </Tabs>
    </>

}

export default StoryForm;