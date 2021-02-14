import {
    Button,
    Textarea,
    VStack,
    Text
} from "@chakra-ui/react";
import { React, useContext, useState } from "react";
import Server from '../server';
import UserContext from '../UserContext';

function GitPush(props) {
    const { project } = useContext(UserContext)
    const [pushInProgress, setPushInProgress] = useState(false)
    const [pushOutput, setPushOutput] = useState(null)

    function setOutput(ok, msg, details) {
        setPushOutput({
            ok: ok,
            msg: msg,
            details: details
        })
        setPushInProgress(false)
    }

    function push() {
        setPushOutput(null)
        setPushInProgress(true)
        Server.postGitPush(project)
            .then(data => setOutput(true, 'All Good!', data))
            .catch(r => {
                const msg = r.response.status == 409 ?
                    'Data on remote server is in conflict. ' +
                    'Try to pull the content before push' :
                    'something went wrong'
                setOutput(false, msg, r.response.data)
            })
    }

    return <VStack spacing={5}>
        <Button size="lg" colorScheme="blue" isLoading={pushInProgress}
            onClick={push}>
            Push
        </Button>
        {pushOutput ? <>
            <Text fontSize="lg" color={pushOutput.ok ? 'green' : 'red'}>
                {pushOutput.msg}
            </Text>
            <Textarea
                value={pushOutput.details}
                size="md"
                resize="Vertical"
                rows="6"
            />
        </> : null}
    </VStack>

}
export default GitPush;