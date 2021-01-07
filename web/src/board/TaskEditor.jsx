import { React, useEffect, useState, useContext } from "react";
import ReactMde from "react-mde";
import ReactMarkdown from "react-markdown";
import "react-mde/lib/styles/css/react-mde-all.css";
import { Box, Button, Center } from "@chakra-ui/react";
import T from "../core/T";
import UserContext from '../UserContext';
import Server from '../server';






function TaskEditor(props) {
  const { project } = useContext(UserContext);
  const { name, readOnly, tags, users } = props;

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
          return tags.map( t => ({ preview: t, value: t}))
        })
    }

  }

  const save = async function* (data) {
    // Promise that waits for "time" milliseconds
    const wait = function (time) {
      return new Promise((a, r) => {
        setTimeout(() => a(), time);
      });
    };

    // Upload "data" to your server
    // Use XMLHttpRequest.send to send a FormData object containing
    // "data"
    // Check this question: https://stackoverflow.com/questions/18055422/how-to-receive-php-image-data-over-copy-n-paste-javascript-with-xmlhttprequest

    await wait(2000);
    // yields the URL that should be inserted in the markdown
    yield "https://picsum.photos/300";
    await wait(2000);

    // returns true meaning that the save was successful
    return true;
  };

  const { task, saveTask } = props
  const [value, setValue] = useState(task && task.description);

  const features = {
    name: "features",
    icon: () => (
      <Button colorScheme="yellow"><T>features</T></Button>
    ),
    execute: opts => {
      opts.textApi.replaceSelection("NICE");
    }
  };

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
    <ReactMde
      value={value}
      onChange={onChange}
      disablePreview={true}
      loadSuggestions={loadSuggestions}
      suggestionTriggerCharacters={['@', '#']}
      paste={{
        saveImage: save
      }}
    /></Box>;
}

export default TaskEditor