import { Box } from "@chakra-ui/react";
import { React, useEffect, useState, useContext } from "react";
import ReactMarkdown from "react-markdown";
import "react-mde/lib/styles/css/react-mde-all.css";


function TaskViewer(props) {
    const { task } = props;

    return task && <Box maxHeight="200px" style={{ overflow: 'auto' }}>
        <div className="mde-preview">
            <div className="mde-preview-content">
                <ReactMarkdown source={task.description} />
            </div>
        </div>
    </Box>

}

export default TaskViewer