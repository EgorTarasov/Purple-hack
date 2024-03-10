import React from 'react';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface ModelSelectProps {
  onSelectChange: (value: string) => void;
}

const ModelSelect: React.FC<ModelSelectProps> = ({ onSelectChange }) => {



  return (
    <Select onValueChange={(value)=>{
      onSelectChange(value);
    }} >
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="LLAMA 2" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Модель</SelectLabel>
          <SelectItem value="llama">LLAMA 2</SelectItem>
          <SelectItem value="fred">FRED T5</SelectItem>
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}

export default ModelSelect;
