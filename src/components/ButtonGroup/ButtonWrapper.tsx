import React from "react";
import isNil from "lodash/isNil";
import Button, { ButtonProps } from "../Button/Button";
import Tooltip from "../Tooltip/Tooltip";
import { HIDE_SHOW_EVENTS } from "../Dialog/consts/dialog-show-hide-event";
import { DialogPositions } from "../../constants/sizes";

export interface ButtonWrapperProps extends ButtonProps {
  tooltipContent?: string;
  tooltipPosition?: typeof DialogPositions[keyof typeof DialogPositions];
  tooltipHideDelay?: number;
  tooltipShowDelay?: number;
  tooltipContainerSelector?: string;
  tooltipMoveBy?: { main?: number; secondary?: number };
}

export const ButtonWrapper = ({
  tooltipContent,
  tooltipPosition,
  tooltipHideDelay,
  tooltipShowDelay,
  tooltipContainerSelector,
  tooltipMoveBy,
  ...otherProps
}: ButtonWrapperProps) => {
  let button = <Button {...otherProps} />;
  if (!isNil(tooltipContent)) {
    button = (
      <Tooltip
        moveBy={tooltipMoveBy}
        position={tooltipPosition}
        hideDelay={tooltipHideDelay}
        showDelay={tooltipShowDelay}
        content={tooltipContent}
        showTrigger={[HIDE_SHOW_EVENTS.MOUSE_ENTER]}
        hideTrigger={[HIDE_SHOW_EVENTS.MOUSE_LEAVE]}
        containerSelector={tooltipContainerSelector}
      >
        {button}
      </Tooltip>
    );
  }

  return button;
};
