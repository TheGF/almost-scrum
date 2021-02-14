import {
    Button,
    Textarea,
    VStack,
    Text
} from "@chakra-ui/react";
import { React, useContext, useState } from "react";
import Server from '../server';
import UserContext from '../UserContext';

function GitPull(props) {
    const { project } = useContext(UserContext)
    const [pullInProgress, setPullInProgress] = useState(false)
    const [pullOutput, setPullOutput] = useState(null)

    function setOutput(ok, msg, details) {
        setPullOutput({
            ok: ok,
            msg: msg,
            details: details
        })
        setPullInProgress(false)
    }

    function pull() {
        setPullOutput(null)
        setPullInProgress(true)
        Server.postGitPull(project)
            .then(data => setOutput(true, 'All Good!', data))
            .catch(r => {
                const msg = 'something went wrong'
                setOutput(false, msg, r.response.data)
            })
    }

    return <VStack spacing={5}>
        <Button size="lg" colorScheme="blue" isLoading={pullInProgress}
            onClick={pull}>
            Pull
        </Button>
        {pullOutput ? <>
            <Text fontSize="lg" color={pullOutput.ok ? 'green' : 'red'}>
                {pullOutput.msg}
            </Text>
            <Textarea
                value={pullOutput.details}
                size="md"
                resize="Vertical"
                rows="6"
            />
        </> : null}
    </VStack>

}
export default GitPull;