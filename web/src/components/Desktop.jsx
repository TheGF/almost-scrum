import React, { useState } from 'react';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import Row from 'react-bootstrap/Row';
import { GiGrapes } from 'react-icons/gi';
import { IoMdLogOut } from 'react-icons/io';
import Project from './project/Project';
import ProjectsList from './project/ProjectsList';

function TopBar() {
    function logout() {
        localStorage.removeItem('username');
        localStorage.removeItem('token');
        window.location.assign(window.location.href);
    }


    return  <Navbar bg="dark" variant="dark">
    <Navbar.Brand href="#home"><GiGrapes/> Almost Scrum 0.1</Navbar.Brand>
    <Nav className="c" >
      <Nav.Link href="#">Git Import</Nav.Link>
      <Nav.Link href="#">Features</Nav.Link>
      <Nav.Link href="#" onClick={logout}>Logout <IoMdLogOut/></Nav.Link>
    </Nav>
    
  </Navbar>
}

function Desktop(props) {
    const [project, setProject] = useState('');

    return <Container fluid>
        <Row>
            <Col md="12">
            <TopBar/>
            </Col>
        </Row>
        <Row>
            <Col style={{flex: 0}}>
                <ProjectsList onSelect={setProject} />
            </Col>
            <Col md>
                <Project username={props.username} project={project} />
            </Col>
        </Row>
    </Container>;
}

export default Desktop;