/**
 * Base Property Panel Component
 * Provides common UI structure for all property panels
 */

import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { Save, Edit3, X, Check } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { EditableNodeProperties, PropertyUpdateCallback, UUID } from '@/types';

interface BasePropertyPanelProps {
  nodeId: UUID;
  nodeType: string;
  title: string;
  description: string;
  children?: React.ReactNode;
  onUpdate: PropertyUpdateCallback;
  className?: string;
}

export const BasePropertyPanel: React.FC<BasePropertyPanelProps> = ({
  nodeId,
  nodeType,
  title,
  description,
  children,
  onUpdate,
  className
}) => {
  const [isEditing, setIsEditing] = useState(false);
  const [editedTitle, setEditedTitle] = useState(title);
  const [editedDescription, setEditedDescription] = useState(description);
  const [hasChanges, setHasChanges] = useState(false);

  const handleStartEdit = () => {
    setIsEditing(true);
    setEditedTitle(title);
    setEditedDescription(description);
    setHasChanges(false);
  };

  const handleCancelEdit = () => {
    setIsEditing(false);
    setEditedTitle(title);
    setEditedDescription(description);
    setHasChanges(false);
  };

  const handleSaveEdit = () => {
    if (hasChanges) {
      onUpdate(nodeId, {
        title: editedTitle,
        description: editedDescription
      });
    }
    setIsEditing(false);
    setHasChanges(false);
  };

  const handleTitleChange = (value: string) => {
    setEditedTitle(value);
    setHasChanges(value !== title || editedDescription !== description);
  };

  const handleDescriptionChange = (value: string) => {
    setEditedDescription(value);
    setHasChanges(editedTitle !== title || value !== description);
  };

  return (
    <Card className={cn("w-full", className)}>
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <div className="flex-1 min-w-0">
            {isEditing ? (
              <div className="space-y-2">
                <div>
                  <Label htmlFor="title" className="text-xs font-medium">
                    Title
                  </Label>
                  <Input
                    id="title"
                    value={editedTitle}
                    onChange={(e) => handleTitleChange(e.target.value)}
                    className="h-8 text-sm font-medium"
                    placeholder="Enter title..."
                  />
                </div>
              </div>
            ) : (
              <CardTitle className="text-sm font-medium leading-tight">
                {title}
              </CardTitle>
            )}
            <Badge variant="outline" className="text-xs mt-1">
              {nodeType}
            </Badge>
          </div>
          
          <div className="flex items-center gap-1 ml-2">
            {isEditing ? (
              <>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-6 w-6 p-0"
                  onClick={handleSaveEdit}
                  disabled={!hasChanges}
                >
                  <Check className="h-3 w-3" />
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-6 w-6 p-0"
                  onClick={handleCancelEdit}
                >
                  <X className="h-3 w-3" />
                </Button>
              </>
            ) : (
              <Button
                variant="ghost"
                size="sm"
                className="h-6 w-6 p-0"
                onClick={handleStartEdit}
              >
                <Edit3 className="h-3 w-3" />
              </Button>
            )}
          </div>
        </div>
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Description Section */}
        <div>
          <Label className="text-xs font-medium text-muted-foreground">
            Description
          </Label>
          {isEditing ? (
            <Textarea
              value={editedDescription}
              onChange={(e) => handleDescriptionChange(e.target.value)}
              className="mt-1 text-xs resize-none"
              rows={3}
              placeholder="Enter description..."
            />
          ) : (
            <p className="mt-1 text-xs leading-relaxed text-muted-foreground">
              {description || "No description provided"}
            </p>
          )}
        </div>

        {children && (
          <>
            <Separator />
            {children}
          </>
        )}
      </CardContent>
    </Card>
  );
};