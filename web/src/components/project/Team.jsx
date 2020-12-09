import React, { useEffect } from 'react';
import Form from 'react-bootstrap/Form';
import Row from 'react-bootstrap/Row';
import Container from 'react-bootstrap/Container';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';

function Team(props) {
    return <Card style={{padding: '1em'}}>
        <Form>
            <Form.Group as={Row}>
                <Form.Label column sm={2}>Team Velocity</Form.Label>
                <Col sm={1}>
                    <Form.Text>2</Form.Text>
                </Col>
                <Col sm={9}>
                    <Form.Control type="range" />
                </Col>
            </Form.Group>
        </Form>
    </Card>;
}

export default Team