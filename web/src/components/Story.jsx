import React, { useEffect, useState } from 'react';
import ListGroup from 'react-bootstrap/ListGroup';
import Button from 'react-bootstrap/Button';
import ButtonGroup from 'react-bootstrap/ButtonGroup';
import ToggleButton from 'react-bootstrap/ToggleButton';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';
import InputGroup from 'react-bootstrap/InputGroup';
import Breadcrumb from 'react-bootstrap/Breadcrumb';
import Badge from 'react-bootstrap/Badge';
import Form from 'react-bootstrap/Form';
import Card from 'react-bootstrap/Card';
import Modal from 'react-bootstrap/Modal';
import axios from 'axios';
import { getConfig, loginWhenUnauthorized } from './axiosUtils';
import { GrNewWindow } from 'react-icons/gr';
import { BiArchive } from 'react-icons/bi';
import { AiOutlineToTop } from 'react-icons/ai';
import { useForm } from 'react-hook-form';
import { emptyStory } from './consts';

const fiboSet = [1, 2, 3, 5, 8, 13, 21];

function Story(props) {
    const { project, store, storyName } = props;
    if (!store || !storyName) return <Badge>Select a story</Badge>;

    const { register, reset, watch } = useForm();
    function fetch() {
        axios.get(`/api/v1/projects/${project}/${store}/${storyName}`, getConfig())
            .then(r => story = reset(r.data))
            .catch(loginWhenUnauthorized);
    }
    useEffect(fetch,[reset]);

    return <>
        <Breadcrumb>
            <Breadcrumb.Item>{store}</Breadcrumb.Item>
            <Breadcrumb.Item>{storyName}</Breadcrumb.Item>
        </Breadcrumb>
        <Form>
            <Form.Group controlId="Description">
                <Form.Label>Description</Form.Label>
                <Form.Control name="description" type="text" placeholder="Description"
                    ref={register} />
            </Form.Group>

            <Form.Group controlId="Description">
            <Form.Label>Points</Form.Label>

            <ButtonGroup toggle>
        {fiboSet.map((number, idx) => (
          <ToggleButton
            key={idx}
            type="radio"
            variant="secondary"
            name="points"
          >
            {number}
          </ToggleButton>
          
        ))}
      </ButtonGroup>
      </Form.Group>
        </Form>


    </>;
}

export default Story;