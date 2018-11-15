import React, { Component } from "react";
import { PageHeader, Panel, PanelGroup, Row, Col } from "react-bootstrap";
import LoaderButton from "../components/LoaderButton";
import config from "../config";
import "./Home.css";


export default class UpdateWidget extends Component {
    constructor(props) {
        super(props);
        var search = this.props.location.search;
        var orderId = parseInt(search.replace("?orderid=", ""))
        if (orderId < 1) {
            orderId = null
        }
        this.state = {
            isLoading: null,
            inventory: [],
            inventoryMask: {},
            orderId: orderId,
            colors: {},
            sizes: {},
            categories: {},
            color: "",
            size: "",
            cageory: ""
        };
    }

    buildFilters = () => {
        const colors = {}; const sizes = {}; const categories = {};
        this.state.inventory.forEach(w => {
            colors[w.Color] = true;
            sizes[w.Size] = true;
            categories[w.Category] = true;
        });
        this.setState({ colors: colors, sizes: sizes, categories: categories });
    }

    updateMask = () => {
        this.state[event.target.id] = event.target.value;
        const mask = this.state.inventoryMask;
        this.state.inventory.forEach(widget => {
            mask[widget.ID] = false
            if (this.state.color && widget.Color != this.state.color) {
                mask[widget.ID] = true
            } else if (this.state.size && widget.Size != this.state.size) {
                mask[widget.ID] = true
            } else if (this.state.category && widget.Category != this.state.category) {
                mask[widget.ID] = true
            }
            if (mask[widget.ID] == true) {
            }

        });
        this.setState({ inventoryMask: mask });
    }



    getWidget = widgetId => {
        return this.state.inventory.find(w => {
            return w.ID == widgetId;
        });
    }

    async componentDidMount() {
        try {
            const resp = await fetch(config.apiGateway.URL + "/widgets");
            const json = await resp.json()
            var data = json.Data || []
            var inventory = []
            data.forEach(widget => {
                widget.CartCount = 0;
                inventory.push(widget);
            });
            this.setState({ inventory });
            // if editing current order
            if (this.state.orderId) {
                const resp = await fetch(config.apiGateway.URL + "/orders/" + this.state.orderId)
                const json = await resp.json()
                var order = json.Data || []
                if (order && order.length > 0) {
                    order[0].LineItems.forEach(li => {
                        console.log(li)
                        console.log(this.state.inventory)
                        const widget = this.getWidget(li.WidgetID)
                        widget.CartCount = li.Quantity
                        widget.Remaining += li.Quantity
                    });
                    this.setState({ inventory });
                }

            }
            this.buildFilters();
        } catch (e) {
            alert(e);
        }
        this.setState({ isLoading: false });
    }

    handleChange = event => {
        const id = parseInt(event.target.id.substr("count-".length))
        var w = this.getWidget(id)
        const val = parseInt(event.target.value)
        w.CartCount = val;
        this.state.inventory.forEach((widget, i) => {
            if (widget.ID === w.ID) {
                w.CartCount = val;
                const tempInventory = this.state.inventory;
                tempInventory[i] = w;
                this.setState({ inventory: tempInventory });
                // TODO : should break
            }
        });
    }

    handleSubmit = async event => {
        try {
            event.preventDefault();
            var cartHasItem = false

            //build order
            const order = {
                id: this.state.orderId,
                lineItems: []
            }
            this.state.inventory.forEach(widget => {
                if (widget.CartCount > 0) {
                    cartHasItem = true
                    order.lineItems.push({
                        widgetId: widget.ID,
                        quantity: widget.CartCount
                    })
                }
            });

            if (!this.state.orderId && !cartHasItem) {
                alert("No widgets selected for purchase!")
                return
            }

            this.setState({ isLoading: true });
            const resp = await fetch(config.apiGateway.URL + "/order", {
                method: "PUT",
                body: JSON.stringify(order)
            })
            const json = await resp.json()
            var data = json.Data

            if (resp.ok) {
                alert("Order created!")
                this.props.history.push("/orders/" + data.ID)
            } else {
                alert("Failed to create order: " + resp.text)
            }
        } catch (e) {
            alert(e);
        }
        this.setState({ isLoading: false });
    }

    renderInventoryList(inventory, mask) {
        return [].concat(inventory).filter(
            (widget) => !mask[widget.ID]
        ).map(
            (widget, i) =>
                widget.Remaining <= 0 ? "" :
                    <Panel key={widget.ID}>
                        <Panel.Heading className={widget.Color}>
                            <Panel.Title componentClass="h3">{widget.Name}</Panel.Title>
                        </Panel.Heading>
                        <Panel.Body>
                            <Row>
                                <Col md={6}>
                                    {"Category: " + widget.Category}
                                </Col>
                                <Col md={1} mdOffset={1}>
                                    {"Stock: " + widget.Remaining}
                                </Col>
                                <Col md={1}>
                                    Quantity:
                        </Col>
                                <Col md={1}>
                                    <input id={"count-" + widget.ID} type="number" min="0" max={widget.Remaining} value={widget.CartCount} onChange={this.handleChange} />
                                </Col>
                            </Row>
                            <Row></Row>
                            {"Color: " + widget.Color}
                            <Row></Row>
                            {"Size:" + widget.Size}
                        </Panel.Body>
                    </Panel>

        );
    }

    renderSelect(type, options, selected) {
        return (
            <select value={selected} onChange={this.updateMask} id={type}>
                <option value=""></option>
                {[].concat(Object.keys(options)).map(
                    option =>
                        <option key={option} value={option}>{option}</option>
                )}
            </select>
        )
    }

    render() {
        return (
            <div className="inventory">
                {this.state.orderId ? <h1>Edit order {this.state.orderId}</h1> : null}
                <form onSubmit={this.handleSubmit}>
                    <LoaderButton
                        block
                        bsStyle="primary"
                        bsSize="large"
                        type="submit"
                        isLoading={this.state.isLoading}
                        text="Place Order"
                        loadingText="creating order..."
                    />
                    <PageHeader>Available Widgets</PageHeader>
                    <PanelGroup id="filters">Filters: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
                            Color: {this.renderSelect("color", this.state.colors, this.state.color)}
                            &nbsp;&nbsp;Category: {this.renderSelect("category", this.state.categories, this.state.category)}
                            &nbsp;&nbsp;Size: {this.renderSelect("size", this.state.sizes, this.state.size)}
                    </PanelGroup>

                    <PanelGroup id="widgets">
                        {!this.state.isLoading && this.renderInventoryList(this.state.inventory, this.state.inventoryMask)}
                    </PanelGroup>
                </form>
            </div>
        );
    }
}

