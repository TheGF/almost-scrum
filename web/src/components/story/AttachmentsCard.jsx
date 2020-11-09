import React, { useState } from 'react';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Button from 'react-bootstrap/Button';
import Card from 'react-bootstrap/Card';
import ListGroup from 'react-bootstrap/ListGroup';
import Modal from 'react-bootstrap/Modal';
import { BiLink, BiTrash, BiUnlink, BiUpload } from 'react-icons/bi';
import { GrDocumentUpdate } from 'react-icons/gr';
import { VscOpenPreview } from 'react-icons/vsc';
import Library from '../library/Library';
import { ButtonToolbar } from 'react-bootstrap';


function AttachmentsCard(props) {
    const { project, form } = props;
    const { getValues, setValue, watch } = form;
    const [showLibrary, setShowLibrary] = useState(false);

    function attachToStory(path) {
        if (!attachments.includes(path)) {
            setValue('attachments', [
                path,
                ...attachments
            ], { shouldValidate: true });
        }
        setShowLibrary(false);
    }

    function LibraryDialog(props) {
        return <Modal show={showLibrary} dialogClassName="modal-80w">
            <Modal.Body>
                <Library project={project} attachToStory={attachToStory} />
            </Modal.Body>
            <Modal.Footer>
                <Button variant="secondary" onClick={_ => setShowLibrary(false)}>
                    Close
                </Button>
            </Modal.Footer>
        </Modal>
    }

    function deleteAttachment(file) {
        setValue('attachments', attachments.filter(
            f => f != file
        ), { shouldValidate: true });
    }

    const attachments = watch('attachments') || [];
    const attachmentsList = attachments.map(f => {
        const parts = f.split('/');
        const name = parts[parts.length - 1];

        return <ListGroup.Item key={f}>
            <Container fluid>
            <Row>
                <Col md="6">
                    <span title={f}>{name}</span>
                </Col>
                <Col md="6" >
                    <ButtonToolbar className="float-right">
                        <Button title="Update the document"><BiUpload /></Button>
                        <Button><VscOpenPreview /></Button>
                        <Button title="Unlink or remove the document" onClick={_ => deleteAttachment(f)}>
                            <BiUnlink /></Button>
                    </ButtonToolbar>
                </Col>
            </Row>
            </Container>
        </ListGroup.Item>;
    }
    );

    return <Card body>
        <LibraryDialog />
        <Button onClick={_ => setShowLibrary(true)}>
            <BiLink /> Add
        </Button>
        <ListGroup>
                {attachmentsList}
        </ListGroup>
    </Card>;
}


export default AttachmentsCard;
