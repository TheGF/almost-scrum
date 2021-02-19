import { Box, Center, HStack } from "@chakra-ui/react";
import { React, useContext, useState } from "react";
import MarkdownEditor from '../core/MarkdownEditor';
import Server from '../server';
import UserContext from '../UserContext';
import Properties from './Properties';



function TaskEditor(props) {
  const { project } = useContext(UserContext);
  const { name, readOnly, tags, users, height } = props;

  function loadSuggestions(text, triggeredBy) {
    if (triggeredBy == '@') {
      return new Promise((accept) => {
        const suggestions = users
          .filter(u => u.includes(text))
          .map(u => ({ preview: u, value: `@${u}`, }))
        accept(suggestions);
      })
    }
    if (triggeredBy == '#') {
      return Server.getSuggestions(project, `%23${text}`, 64)
        .then(tags => {
          return tags.map(t => ({ preview: t, value: t }))
        })
    }

  }

  const { task, saveTask } = props
  const [value, setValue] = useState(task && task.description);

  function onChange(value) {
    setValue(value);
    if (task) {
      task.description = value;
      saveTask(task);
    }
  }
  const editMessage = readOnly ? <Center h="3em">
    Change owner if you want to edit the content
</Center> : null

  return readOnly ? editMessage : <Box>
    <HStack>
      <Box w="90%" h={height}>
      <MarkdownEditor
        value={value}
        height={height}
        onChange={onChange}
        disablePreview={true}
        loadSuggestions={loadSuggestions}
        suggestionTriggerCharacters={['@', '#']}
        imageFolder="/.inline-images"
      />
      </Box>
      <Box minW="240px" h={height-20}>
      <Properties task={task} saveTask={saveTask} users={users} height={height}
        height={height} readOnly={readOnly} />
      </Box>
    </HStack>
  </Box>;
}

export default TaskEditor