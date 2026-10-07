import Window from '../../../src/window/Window';
import Document from '../../../src/nodes/document/Document';
import HTMLSelectElement from '../../../src/nodes/html-select-element/HTMLSelectElement';
import HTMLOptionElement from '../../../src/nodes/html-option-element/HTMLOptionElement';
import DOMException from '../../../src/exception/DOMException';

describe('HTMLSelectElement', () => {
	let window: Window;
	let document: Document;
	let element: HTMLSelectElement;

	beforeEach(() => {
		window = new Window();
		document = window.document;
		element = <HTMLSelectElement>document.createElement('select');
	});

	describe('Object.prototype.toString', () => {
		it('Returns `[object HTMLSelectElement]`', () => {
			expect(Object.prototype.toString.call(element)).toBe('[object HTMLSelectElement]');
		});
	});

	describe('get value()', () => {
		it('Returns the value of the selected option.', () => {
			element.innerHTML =
				'<option value="value1">Value 1</option><option value="value2" selected>Value 2</option>';
			expect(element.value).toBe('value2');
		});

		it('Returns an empty string when no option is selected.', () => {
			element.innerHTML = '<option value="value1">Value 1</option>';
			expect(element.value).toBe('');
		});
	});

	describe('set value()', () => {
		it('Selects the option with a matching value.', () => {
			element.innerHTML =
				'<option value="value1">Value 1</option><option value="value2">Value 2</option>';
			element.value = 'value2';
			expect(element.selectedIndex).toBe(1);
			expect((<HTMLOptionElement>element.options[1]).selected).toBe(true);
			expect(element.options[1].getAttribute('selected')).toBe('');
			expect((<HTMLOptionElement>element.options[0]).selected).toBe(false);
		});

		it('Sets selectedIndex to -1 when no option matches.', () => {
			element.innerHTML = '<option value="value1" selected>Value 1</option>';
			element.value = 'invalid';
			expect(element.selectedIndex).toBe(-1);
			expect(element.value).toBe('');
		});

		it('Trims and removes new lines.', () => {
			element.innerHTML = '<option value="value1">Value 1</option>';
			element.value = '  \nvalue1\n  ';
			expect(element.selectedIndex).toBe(0);
		});
	});

	describe('get selectedIndex()', () => {
		it('Returns -1 when no option is selected.', () => {
			element.innerHTML = '<option value="value1">Value 1</option>';
			expect(element.selectedIndex).toBe(-1);
		});

		it('Returns the index of the option with a "selected" attribute in the initial HTML.', () => {
			element.innerHTML =
				'<option value="value1">Value 1</option><option value="value2" selected>Value 2</option>';
			expect(element.selectedIndex).toBe(1);
		});
	});

	describe('set selectedIndex()', () => {
		it('Updates the "selected" attribute on the option nodes.', () => {
			element.innerHTML =
				'<option value="value1" selected>Value 1</option><option value="value2">Value 2</option>';
			element.selectedIndex = 1;
			expect(element.options[0].getAttribute('selected')).toBe(null);
			expect(element.options[1].getAttribute('selected')).toBe('');
		});

		it('Unselects all options when set to -1.', () => {
			element.innerHTML =
				'<option value="value1" selected>Value 1</option><option value="value2">Value 2</option>';
			element.selectedIndex = -1;
			expect(element.selectedIndex).toBe(-1);
			expect(element.options[0].getAttribute('selected')).toBe(null);
			expect(element.options[1].getAttribute('selected')).toBe(null);
		});

		it('Throws an error when the index is out of bounds.', () => {
			element.innerHTML = '<option value="value1">Value 1</option>';
			expect(() => {
				element.selectedIndex = 1;
			}).toThrow(DOMException);
			expect(() => {
				element.selectedIndex = -2;
			}).toThrow(DOMException);
		});
	});

	describe('get options()', () => {
		it('Reflects options added to and removed from the DOM.', () => {
			expect(element.options.length).toBe(0);
			element.innerHTML =
				'<option value="value1">Value 1</option><option value="value2">Value 2</option>';
			expect(element.options.length).toBe(2);
			expect((<HTMLOptionElement>element.options[0]).value).toBe('value1');
			expect((<HTMLOptionElement>element.options[1]).value).toBe('value2');
			element.removeChild(element.options[0]);
			expect(element.options.length).toBe(1);
			expect((<HTMLOptionElement>element.options[0]).value).toBe('value2');
		});

		it('Includes options nested in an optgroup.', () => {
			element.innerHTML =
				'<optgroup label="group"><option value="value1">Value 1</option></optgroup>';
			expect(element.options.length).toBe(1);
			expect((<HTMLOptionElement>element.options[0]).value).toBe('value1');
		});
	});

	for (const property of ['disabled', 'autofocus', 'required', 'multiple']) {
		describe(`get ${property}()`, () => {
			it('Returns attribute value.', () => {
				expect(element[property]).toBe(false);
				element.setAttribute(property, '');
				expect(element[property]).toBe(true);
			});
		});

		describe(`set ${property}()`, () => {
			it('Sets attribute value.', () => {
				element[property] = true;
				expect(element.getAttribute(property)).toBe('');
			});
		});
	}

	describe(`get name()`, () => {
		it('Returns attribute value.', () => {
			expect(element.name).toBe('');
			element.setAttribute('name', 'value');
			expect(element.name).toBe('value');
		});
	});

	describe(`set name()`, () => {
		it('Sets attribute value.', () => {
			element.name = 'value';
			expect(element.getAttribute('name')).toBe('value');
		});
	});
});
