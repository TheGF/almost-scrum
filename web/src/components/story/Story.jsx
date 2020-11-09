import React, { useEffect, useState } from 'react';
import Badge from 'react-bootstrap/Badge';
import { useForm } from 'react-hook-form';
import Server from '../Server';
import StoryForm from './StoryForm';
import Utils from '../utils';


function Story(props) {
    const { project, store, story } = props;
    if (!store || !story) return <Badge>Select a story</Badge>;

    function lazySave(values) {
        const key = `/${project}/${store}/${story}`;
        setPendingWrite(true);
        Utils.lazyCall(key, _ =>
            Server.saveStory(project, store, story, values)
                .then(_ => setPendingWrite(false)));
        return {
            values: values, errors: {}
        }
    }
    const form = useForm({ mode: 'onChange', resolver: lazySave });

    function fetch() {
        Server.getStory(project, store, story)
            .then(form.reset)
            .then(Utils.sendPendingCalls);
    }
    useEffect(fetch, [project, store, story]);


    const [pendingWrite, setPendingWrite] = useState(0);

    return <StoryForm {...props} form={form} pendingWrite={pendingWrite} />
}

export default Story;