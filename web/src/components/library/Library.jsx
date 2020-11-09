import React, { useEffect, useState } from 'react';
import { Breadcrumb, ButtonToolbar } from 'react-bootstrap';
import Button from 'react-bootstrap/Button';
import ButtonGroup from 'react-bootstrap/ButtonGroup';
import Card from 'react-bootstrap/Card';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import ListGroup from 'react-bootstrap/ListGroup';
import Row from 'react-bootstrap/Row';
import { BiHomeAlt, BiLink, BiRefresh, BiTrash } from 'react-icons/bi';

import Server from '../Server';
import Utils from '../utils';
import Form from 'react-bootstrap/Form';

function Library(props) {
    const { project, attachToStory } = props;
    const [path, setPath] = useState(props.path || localStorage.libraryPath || '');
    const [files, setFiles] = useState([]);

    function fetch() {
        Server.libraryList(project, path)
            .then(setFiles);
    }
    useEffect(fetch, [path])

    function fileSelect(file) {
        if (file.dir) {
            setPath(`${path}/${file.name}`);
        } else {
            Server.libraryDownload(project, `${path}/${file.name}`);
        }
    }

    function fileUpload(evt) {
        const file = evt.target.files[0];
        file && Server.libraryPost(project, path, file)
            .then(setFiles);
    }

    function fileDelete(path) {
        Server.libraryDelete(project, path);
        fetch();
    }

    function fileView(file) {
        const { name, dir, mime, mtime, creator } = file;
        const time = Utils.getFriendlyDate(mtime * 1000);

        return <Row>
            <Col md="5">
                <a href="#" onClick={_ => fileSelect(file)}>
                    {Utils.fileIcon(dir, mime)}
                &nbsp;
                {name}
                </a>
            </Col>
            <Col md="2">
                {creator && `${creator} `}
            </Col>
            <Col md="3">
                {time}
            </Col>
            <Col md="2" className="float-right">
                <ButtonGroup>
                    {!dir && attachToStory ?
                        <Button onClick={_ => attachToStory(`${path}/${name}`)}>
                            <BiLink /> Add To Story
                        </Button>
                        : null
                    }
                    <Button onClick={_ => fileDelete(`${path}/${name}`)}>
                        <BiTrash />Delete
                    </Button>
                </ButtonGroup>
            </Col>
        </Row>;
    }

    const filesView = files.map((f, idx) =>
        <ListGroup.Item key={idx}>
            {fileView(f)}
        </ListGroup.Item>
    );

    const parts = path.split('/').reduce((acc, p) => {
        const last = acc.length ? acc[acc.length - 1][1] : '';
        acc.push([p || <BiHomeAlt />, `${last}/${p}`]);
        return acc;
    }, []);
    const crumbs = parts.map(p =>
        <Breadcrumb.Item key={p} href="#" onClick={_ => setPath(p[1])}>
            {p[0]}
        </Breadcrumb.Item>
    );

    return <Card>
        <Form>
            <ButtonToolbar>
                <Button onClick={fetch}><BiRefresh />Refresh</Button>
                <Form.File
                    label="Add to the Library"
                    id="custom-file"
                    custom
                    onChange={fileUpload}
                />
            </ButtonToolbar>
            <Breadcrumb>{crumbs}</Breadcrumb>
        </Form>
        <Container fluid>
            <Row>
                <Col md="6">Name</Col>
                <Col md="2">Owner</Col>
                <Col md="2">Last Change</Col>
                <Col md="2">Actions</Col>
            </Row>
            {filesView}
            <br />
        </Container>
    </Card>
}

export default Library;