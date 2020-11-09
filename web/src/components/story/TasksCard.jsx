import React, { useState } from 'react';
import Button from 'react-bootstrap/Button';
import ButtonGroup from 'react-bootstrap/ButtonGroup';
import Card from 'react-bootstrap/Card';
import Form from 'react-bootstrap/Form';
import ListGroup from 'react-bootstrap/ListGroup';
import Modal from 'react-bootstrap/Modal';
import ToggleButton from 'react-bootstrap/ToggleButton';
import { useForm } from 'react-hook-form';

function TasksCard(props) {
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
        return <ListGroup.Item key={idx} onClick={
            _ => setEditTaskId(idx)} style={{ cursor: 'pointer' }
            }>
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

export default TasksCard;