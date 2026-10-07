import DOMException from '../../exception/DOMException';
import HTMLCollection from '../element/HTMLCollection';
import HTMLOptGroupElement from '../html-opt-group-element/HTMLOptGroupElement';
import HTMLOptionElement from './HTMLOptionElement';
import IHTMLSelectElement from '../html-select-element/IHTMLSelectElement';
import IHTMLOptionsCollection from './IHTMLOptionsCollection';

/**
 * HTML Options Collection.
 *
 * Reference:
 * https://developer.mozilla.org/en-US/docs/Web/API/HTMLOptionsCollection.
 */
export default class HTMLOptionsCollection
	extends HTMLCollection
	implements IHTMLOptionsCollection
{
	public _selectElement: IHTMLSelectElement;

	/**
	 * Constructor.
	 *
	 * @param selectElement Select element.
	 */
	constructor(selectElement?: IHTMLSelectElement) {
		super();
		this._selectElement = selectElement || null;
	}

	/**
	 * Returns selectedIndex.
	 *
	 * @returns SelectedIndex.
	 */
	public get selectedIndex(): number {
		return this.findIndex((option) => (<HTMLOptionElement>option).selected);
	}

	/**
	 * Sets selectedIndex.
	 *
	 * @param selectedIndex SelectedIndex.
	 */
	public set selectedIndex(selectedIndex: number) {
		for (let i = 0; i < this.length; i++) {
			(<HTMLOptionElement>this[i]).selected = i === selectedIndex;
		}
	}

	/**
	 * Returns item by index.
	 *
	 * @param index Index.
	 */
	public item(index: number): HTMLOptionElement | HTMLOptGroupElement {
		return this[index];
	}

	/**
	 *
	 * @param element
	 * @param before
	 */
	public add(
		element: HTMLOptionElement | HTMLOptGroupElement,
		before?: number | HTMLOptionElement | HTMLOptGroupElement
	): void {
		if (!before && before !== 0) {
			this.push(element);
			if (this._selectElement) {
				this._selectElement.appendChild(element);
			}
			return;
		}

		if (!Number.isNaN(Number(before))) {
			if (before < 0) {
				return;
			}

			this.splice(<number>before, 0, element);
			if (this._selectElement) {
				const beforeElement = this[<number>before + 1];
				if (beforeElement) {
					beforeElement.parentNode.insertBefore(element, beforeElement);
				} else {
					this._selectElement.appendChild(element);
				}
			}
			return;
		}

		const idx = this.findIndex((element) => element === before);
		if (idx === -1) {
			throw new DOMException(
				"Failed to execute 'add' on 'DOMException': The node before which the new node is to be inserted is not a child of this node."
			);
		}

		this.splice(idx, 0, element);
		if (this._selectElement) {
			(<HTMLOptionElement | HTMLOptGroupElement>before).parentNode.insertBefore(
				element,
				<HTMLOptionElement | HTMLOptGroupElement>before
			);
		}
	}

	/**
	 * Removes indexed element from collection.
	 *
	 * @param index Index.
	 */
	public remove(index: number): void {
		const selectedIndex = this.selectedIndex;
		const element = this[index];
		this.splice(index, 1);
		if (element && element.parentNode) {
			element.parentNode.removeChild(element);
		}
		if (index === selectedIndex) {
			this.selectedIndex = this.length ? 0 : -1;
		}
	}
}
