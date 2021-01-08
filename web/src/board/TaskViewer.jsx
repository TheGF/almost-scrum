import { Badge, Box } from "@chakra-ui/react";
import { React } from "react";
import ReactMarkdown from "react-markdown";
import "react-mde/lib/styles/css/react-mde-all.css";
import gfm from 'remark-gfm';


function TaskViewer(props) {
    const { task, searchKeys } = props;

    function highlightKeys(text) {
        let keyLocations = []
        const out = []
        const lower = text.toLowerCase()
        for (const key of searchKeys) {
            let start = 0
            let end = 0
            const lowerKey = key.toLowerCase()
            do {
                start = lower.indexOf(lowerKey, end)
                if (start >= 0) {
                    end = start+key.length
                    keyLocations.push([start, end])
                }
            } while (start != -1)
        }

        keyLocations = keyLocations.sort( (a,b) => a[0]-b[0])
        let cursor = 0 
        for (const [start, end] of keyLocations) {
            if (cursor > start) continue

            out.push(text.substring(cursor, start))
            out.push(<Badge key={out.length} colorScheme="green">
                {text.substring(start, end)}
            </Badge>)
            cursor = end
        }
        out.push(text.substring(cursor))
        return out
    }

    const renderers = { 
        text: (props) => {
            const text = highlightKeys(props.value) 
            return text
        }
    }

    return task && <Box maxHeight="200px" style={{ overflow: 'auto' }}>
        <div className="mde-preview">
            <div className="mde-preview-content">
                <ReactMarkdown plugins={[gfm]} renderers={renderers} 
                source={task.description} />
            </div>
        </div>
    </Box>

}

export default TaskViewer