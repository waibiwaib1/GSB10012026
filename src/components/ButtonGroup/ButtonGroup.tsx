import React, { useRef, forwardRef, useCallback, useMemo, useEffect, useState } from "react";
import cx from "classnames";
import Button from "../../components/Button/Button";
import { ButtonType, Size } from "../../components/Button/ButtonConstants";
import usePrevious from "../../hooks/usePrevious";
import useMergeRefs from "../../hooks/useMergeRefs";
import { backwardCompatibilityForProperties } from "../../helpers/backwardCompatibilityForProperties";
import { baseClassName } from "./ButtonGroupConstants";
import { ButtonWrapper } from "./ButtonWrapper";
import VibeComponentProps from "../../types/VibeComponentProps";
import VibeComponent from "../../types/VibeComponent";
import { SIZES, DialogPositions } from "../../constants/sizes";
import "./ButtonGroup.scss";

export type ButtonGroupValue = string | number;

export interface ButtonGroupOption {
  value: ButtonGroupValue;
  text: string;
  subText?: string;
  disabled?: boolean;
  icon?: string | React.FunctionComponent | null;
  leftIcon?: string | React.FunctionComponent | null;
  ariaLabel?: string;
  tooltipContent?: string;
}

interface ButtonGroupProps extends VibeComponentProps {
  /**
   * Backward compatibility for props naming - please use className instead
   * @deprecated
   */
  componentClassName?: string;
  value?: ButtonGroupValue;
  onSelect?: (value: ButtonGroupValue, name: string) => void;
  name?: string;
  disabled?: boolean;
  options: Array<ButtonGroupOption>;
  size?: Size;
  kind?: ButtonType;
  groupAriaLabel?: string;
  tooltipPosition?: typeof DialogPositions[keyof typeof DialogPositions];
  tooltipHideDelay?: number;
  tooltipShowDelay?: number;
  tooltipContainerSelector?: string;
  tooltipMoveBy?: { main?: number; secondary?: number };
}

const ButtonGroup: VibeComponent<ButtonGroupProps> & {
  sizes?: typeof SIZES;
  kinds?: typeof ButtonType;
} = forwardRef(
  (
    {
      className,
      // Backward compatibility for props naming
      componentClassName,
      options,
      name = "",
      disabled = false,
      value = "",
      onSelect,
      size = SIZES.SMALL,
      kind = ButtonType.SECONDARY,
      groupAriaLabel = "",
      tooltipPosition,
      tooltipHideDelay,
      tooltipShowDelay,
      tooltipContainerSelector,
      tooltipMoveBy
    },
    ref
  ) => {
    const overrideClassName = backwardCompatibilityForProperties([className, componentClassName]);
    const inputRef = useRef();
    const [valueState, setValueState] = useState(value);
    const prevValue = usePrevious(value);
    const mergedRef = useMergeRefs({ refs: [ref, inputRef] });

    const onClick = useCallback(
      (option: ButtonGroupOption) => {
        const isDisabled = disabled || option.disabled;
        if (!isDisabled) {
          setValueState(option.value);
          if (onSelect) {
            onSelect(option.value, name);
          }
        }
      },
      [onSelect, disabled, name]
    );

    const selectedOption = useMemo(() => {
      return options.find(option => option.value === valueState);
    }, [options, valueState]);

    const Buttons = useMemo(() => {
      return options.map((option, index) => {
        const isSelected = option.value === valueState;
        return (
          <ButtonWrapper
            key={option.value}
            size={size}
            onClick={() => onClick(option)}
            rightIcon={option.icon}
            leftIcon={option.leftIcon}
            active={isSelected}
            rightFlat={index !== options.length - 1}
            leftFlat={index !== 0}
            kind={Button.kinds.TERTIARY}
            preventClickAnimation
            ariaLabel={option.ariaLabel}
            tooltipContent={option.tooltipContent}
            tooltipPosition={tooltipPosition}
            tooltipHideDelay={tooltipHideDelay}
            tooltipShowDelay={tooltipShowDelay}
            tooltipContainerSelector={tooltipContainerSelector}
            tooltipMoveBy={tooltipMoveBy}
            className={cx(`${baseClassName}__option-text`, {
              selected: isSelected,
              disabled,
              "button-disabled": option.disabled
            })}
          >
            {option.text}
          </ButtonWrapper>
        );
      });
    }, [
      options,
      disabled,
      onClick,
      size,
      valueState,
      tooltipPosition,
      tooltipHideDelay,
      tooltipShowDelay,
      tooltipContainerSelector,
      tooltipMoveBy
    ]);

    // Effects
    useEffect(() => {
      // Update value if changed from props
      if (value !== prevValue && value !== valueState) {
        setValueState(value);
      }
    }, [value, prevValue, valueState, setValueState]);

    return (
      <div
        className={cx(baseClassName, overrideClassName, `${baseClassName}--kind-${kind}`, { disabled })}
        ref={mergedRef}
      >
        <div
          role="group"
          aria-label={groupAriaLabel}
          className={cx(`${baseClassName}__buttons-container`)}
          aria-disabled={disabled}
        >
          {Buttons}
        </div>
        {selectedOption && selectedOption.subText && (
          <div className={`${baseClassName}__sub-text-container`}>{selectedOption.subText}</div>
        )}
      </div>
    );
  }
);

Object.assign(ButtonGroup, {
  sizes: Button.sizes,
  kinds: Button.kinds
});

export default ButtonGroup;
