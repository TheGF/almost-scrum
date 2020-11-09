import React, { useEffect } from 'react';
import Badge from 'react-bootstrap/Badge';
import Tab from 'react-bootstrap/Tab';
import Tabs from 'react-bootstrap/Tabs';
import { GiSave } from 'react-icons/gi';
import AttachmentsCard from './AttachmentsCard';
import MainCard from './MainCard';
import TasksCard from './TasksCard';

function StoryForm(props) {
    const { project, store, story, form, pendingWrite } = props;
    const { register } = form;

    useEffect(() => register('points'), [register]);
    useEffect(() => register('tasks'), [register]);
    useEffect(() => register('attachments'), [register]);

    return <>
        <span className="float-right">
            {pendingWrite ? <GiSave /> : null}
            <Badge>{project}/{store}/{story}</Badge>
        </span>
        <Tabs defaultActiveKey="main" >
            <Tab title="Main" eventKey="main">
                <MainCard {...props} />
            </Tab>
            <Tab title="Tasks" eventKey="tasks">
                <TasksCard {...props} />
            </Tab>
            <Tab title="Time" eventKey="time" >

            </Tab>
            <Tab title="Attachments" eventKey="attachments" >
                <AttachmentsCard {...props}/>
            </Tab>
        </Tabs>
    </>

}

export default StoryForm;