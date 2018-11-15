import React, { Component } from "react";
import { FormGroup, FormControl, ControlLabel } from "react-bootstrap";
import LoaderButton from "../components/LoaderButton";
import config from "../config";
import "./Admin.css";


export default class UpdateWidget extends Component {
    constructor(props) {
      super(props);
  
      this.state = {
        isLoading: null,
        name: null,
        color: null,
        category: null,
        size: null,
        remaining: 0
      };
    }
  
    validateForm() {
      return this.state.name  && this.state.color && this.state.category && this.state.size && this.state.remaining;
    }
  
    handleChange = event => {
      this.setState({
        [event.target.id]: event.target.type === 'number' ?  parseInt(event.target.value) : event.target.value
      });
    }

    handleSubmit = async event => {
      event.preventDefault();
  
      if (!this.state.name || !this.state.color || !this.state.category || !this.state.size || !this.state.remaining ){
        alert(`Please make sure all fields are set`);
        return;
      }
  
      this.setState({ isLoading: true });
      // fix
      fetch(config.apiGateway.URL+"/widget", {
          method: "PUT",
          body: JSON.stringify(this.state)
      })
      .then(response => {
          this.setState({
              isLoading: false
          })
          if (response.ok) {
            alert("Widget entry updated!")
            this.props.history.push("/")
          } else {
            alert("Failed to update widget entry: " + response.text)
          }
      })
    }
  
    render() {
      return (
        <div className="UpdateProduct">
          <form onSubmit={this.handleSubmit}>
            <FormGroup controlId="name">
            <ControlLabel>Widget Name</ControlLabel>
              <FormControl
                onChange={this.handleChange}
                componentClass="input"
              />
            </FormGroup>
            <FormGroup controlId="category">
            <ControlLabel>Widget Category</ControlLabel>
              <FormControl
                onChange={this.handleChange}
                componentClass="input"
              />
            </FormGroup>
            <FormGroup controlId="size">
            <ControlLabel>Widget Size</ControlLabel>
              <FormControl
                onChange={this.handleChange}
                componentClass="input"
              />
            </FormGroup>
            <FormGroup controlId="color">
            <ControlLabel>Widget Color</ControlLabel>
            <FormControl
              onChange={this.handleChange}
              componentClass="input"
            />
          </FormGroup>
          <FormGroup controlId="remaining">
            <ControlLabel>Stock</ControlLabel>
              <FormControl
                type="number"
                onChange={this.handleChange}
                componentClass="input"
              />
            </FormGroup>
            <LoaderButton
              block
              bsStyle="primary"
              bsSize="large"
              disabled={!this.validateForm()}
              type="submit"
              isLoading={this.state.isLoading}
              text="Update Inventory"
              loadingText="Updating…"
            />
          </form>
        </div>
      );
    }
  }