import React, { useEffect, useState } from 'react';
import Button from 'react-bootstrap/Button';
import ButtonGroup from 'react-bootstrap/ButtonGroup';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';
import Card from 'react-bootstrap/Card';
import Form from 'react-bootstrap/Form';
import Modal from 'react-bootstrap/Modal';
import ToggleButton from 'react-bootstrap/ToggleButton';
import { BiArchive } from 'react-icons/bi';
import { CgMoveTask } from 'react-icons/cg';
import { pointsUnitOptions, statusOptions } from '../consts';
import Server from '../Server';

function MainCard(props) {
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
        Server.getUsers(project)
            .then(setUsers);
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


export default MainCard;