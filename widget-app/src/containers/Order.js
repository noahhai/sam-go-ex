import React, { Component } from "react";
import { PageHeader, Panel, PanelGroup, Row, Col } from "react-bootstrap";
import { LinkContainer, Button } from "react-router-bootstrap";
import LoaderButton from "../components/LoaderButton";
import config from "../config";
import "./Home.css";


export default class OrderWidget extends Component {
    constructor(props) {
        super(props);
        var path = this.props.location.pathname;
        var orderId = parseInt(path.replace("/orders/", ""))
        if (orderId < 1) {
            orderId = null
        }
        this.state = {
            isLoading: null,
            orders: [],
            orderId: orderId
        };

    }

    gotoEditOrder = orderId => { }

    async componentDidMount() {
        try {
            const path = this.state.orderId ? "/orders/" + this.state.orderId : "/orders";
            const resp = await fetch(config.apiGateway.URL + path);
            const json = await resp.json()
            var orders = json.Data || []
            console.log(orders)
            this.setState({ orders: orders });
        } catch (e) {
            alert(e);
        }
        this.setState({ isLoading: false });
    }

    renderInventoryList(inventory) {
        return [].concat(inventory).map(
            (order, i) =>
                <div>

                    <Panel key={order.ID}>
                        <Panel.Heading>
                            <Panel.Title componentClass="h3">Order #{order.ID} &nbsp;&nbsp;&nbsp;
                                <input type="button" value="Edit order" onClick={() => { this.props.history.push("/?orderid=" + order.ID) }} />
                            </Panel.Title>
                        </Panel.Heading>
                        <Panel.Body>
                            {[].concat(order.LineItems).map(
                                function (li, i) {
                                    return (
                                        <div>
                                            Item #{i + 1} => Quantity: {li.Quantity}, Category: {li.Widget.Category}, Name:
                                {li.Widget.Name}, Color: {li.Widget.Color}, Size: {li.Widget.Size}
                                            <Row></Row>
                                        </div>
                                    );
                                })}
                        </Panel.Body>
                    </Panel>
                </div>
        );
    }

    render() {
        return (
            <div className="inventory">
                <form onSubmit={this.handleSubmit}>
                    <PageHeader>Orders</PageHeader>
                    <PanelGroup id="orders">
                        {!this.state.isLoading && this.renderInventoryList(this.state.orders)}
                    </PanelGroup>
                </form>
            </div>
        );
    }
}